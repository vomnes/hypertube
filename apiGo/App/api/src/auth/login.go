package auth

import (
	"log"
	"net/http"
	"time"

	"../../../lib"
	coltype "../../../mongodb/collections"
	"../../../mongodb/query"
	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func checkUserSecret(inputData loginData, db *mgo.Database) (coltype.User, int, string) {
	var u coltype.User
	u, err := query.FindUser(bson.M{"username": inputData.Username, "account_type.level": 1}, db)
	if err != nil && err.Error() != "Not Found" {
		log.Println(lib.PrettyError("[DB REQUEST] Failed to get user data " + err.Error()))
		return coltype.User{}, 500, "User data collection failed"
	}
	if err != nil && (err.Error() == "Not Found" || u == (coltype.User{})) {
		return coltype.User{}, 403, "Username or password incorrect"
	}
	// Comparing the password with the hashed password from the body
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputData.Password))
	if err != nil {
		return coltype.User{}, 403, "Username or password incorrect"
	}
	return u, 0, ""
}

func Login(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.MongoDB).(*mgo.Database)
	if !ok {
		lib.RespondWithErrorHTTP(w, http.StatusInternalServerError, "Problem with database connection")
		return
	}
	var inputData loginData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	if inputData.Username == "" || inputData.Password == "" {
		lib.RespondWithErrorHTTP(w, 403, "Cannot have an empty field")
		return
	}
	u, errCode, errContent := checkUserSecret(inputData, db)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	// Handle login - Return JWT
	settingsJWT := lib.DataJWT{
		Duration:       time.Hour * time.Duration(24*31),
		ISS:            "hypertube.com",
		Sub:            "",
		UserID:         u.ID,
		Username:       u.Username,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		ProfilePicture: u.ProfilePicture,
	}
	jwt, err := lib.GenerateJWT(settingsJWT)
	if err != nil {
		log.Println(lib.PrettyError("Login - " + inputData.Username + " : " + err.Error()))
		lib.RespondWithErrorHTTP(w, http.StatusInternalServerError, "Failed to generate JSON Web Token")
		return
	}
	lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"token": jwt,
	})
}
