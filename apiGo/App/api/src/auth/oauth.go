package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"../../../lib"
	coltypes "../../../mongodb/collections"
	"../../../mongodb/query"
	api42 "./goth/42"
	apifb "./goth/fb"

	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/gplus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func InitOAuth2() {
	goth.UseProviders(
		gplus.New(os.Getenv("GPLUS_KEY"), os.Getenv("GPLUS_SECRET"), os.Getenv("API_DOMAIN_NAME")+"/api/auth/gplus/callback"),
		api42.New(os.Getenv("API42_KEY"), os.Getenv("API42_SECRET"), os.Getenv("API_DOMAIN_NAME")+"/api/auth/42/callback"),
		apifb.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), os.Getenv("API_DOMAIN_NAME")+"/api/auth/facebook/callback"),
	)
}

func oAuthentication(w http.ResponseWriter, r *http.Request, user goth.User) {
	db, ok := r.Context().Value(lib.MongoDB).(*mgo.Database)
	if !ok {
		log.Println(lib.PrettyError("OAuth - Database Connection Failed"))
		lib.RespondWithErrorHTTP(w, http.StatusInternalServerError, "Database connection failed")
		return
	}
	username := user.Provider + "_" + user.FirstName + "_" + user.LastName
	userID := lib.GetRandomString(42)
	search := bson.M{"username": username, "account_type.type": user.Provider}
	dbUser, err := query.FindUser(search, db)
	if err != nil {
		if err.Error() == "Not Found" {
			// Register with OAuth2 data
			profile := coltypes.User{
				ID:             userID,
				Username:       username,
				FirstName:      user.FirstName,
				LastName:       user.LastName,
				ProfilePicture: user.AvatarURL,
				Email:          user.Email,
				AccountType: coltypes.AccountType{
					Level: 2,
					Type:  user.Provider,
				},
				Locale: "en",
			}
			if err := query.InsertUser(profile, db); err != nil {
				log.Println(lib.PrettyError("OAuth - " + username + " : " + err.Error()))
				lib.RespondWithErrorHTTP(w, http.StatusInternalServerError, "Failed to insert user in database")
				return
			}
		} else {
			log.Println(lib.PrettyError("OAuth - " + username + " : " + err.Error()))
			lib.RespondWithErrorHTTP(w, http.StatusInternalServerError, "Failed to collect data in database")
			return
		}
	}
	// dbUser fields are empty in case of registration
	if dbUser.ID == "" {
		dbUser.ID = userID
	}
	if dbUser.FirstName == "" {
		dbUser.FirstName = user.FirstName
	}
	if dbUser.LastName == "" {
		dbUser.LastName = user.LastName
	}
	if dbUser.ProfilePicture == "" {
		dbUser.ProfilePicture = user.AvatarURL
	}
	// Handle login - Return JWT
	settingsJWT := lib.DataJWT{
		Duration:       time.Hour * time.Duration(24*31),
		ISS:            "hypertube.com",
		Sub:            "",
		UserID:         dbUser.ID,
		Username:       username,
		FirstName:      dbUser.FirstName,
		LastName:       dbUser.LastName,
		ProfilePicture: dbUser.ProfilePicture,
		Oauth: lib.Oauth{
			Token:    user.AccessToken,
			Provider: user.Provider,
		},
	}
	jwt, err := lib.GenerateJWT(settingsJWT)
	if err != nil {
		log.Println(lib.PrettyError("OAuth - " + username + " : " + err.Error()))
		lib.RespondWithErrorHTTP(w, http.StatusInternalServerError, "Failed to generate JSON Web Token")
		return
	}
	http.Redirect(w, r, os.Getenv("FRONT_DOMAIN_NAME")+"/gallery?token="+jwt, http.StatusSeeOther)
}

func OAuthProviderCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Println(lib.PrettyError("OAuth - Failed to complete user auth Goth : " + err.Error()))
		http.Redirect(w, r, os.Getenv("FRONT_DOMAIN_NAME")+"/?error=oauth_failed", http.StatusSeeOther)
		return
	}
	oAuthentication(w, r, user)
	return
}

func OAuthProvider(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	providers := goth.GetProviders()
	// Check if called provider is handled
	if providers[vars["provider"]] == nil {
		lib.RespondWithErrorHTTP(w, http.StatusNotFound, "OAuth provider not handled")
		return
	}
	// try to get the user without re-authenticating
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		fmt.Println(gothUser)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}
