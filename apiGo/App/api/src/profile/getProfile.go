package profile

import (
	"net/http"

	"../../../lib"
	coltype "../../../mongodb/collections"
	"../../../mongodb/query"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.MongoDB).(*mgo.Database)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	username, _ := r.Context().Value(lib.Username).(string)
	userID, _ := r.Context().Value(lib.UserID).(string)
	// Get user data
	var user coltype.User
	user, err := query.FindUser(bson.M{"_id": userID, "username": username}, db)
	if err != nil {
		if err.Error() == "Not Found" {
			lib.RespondWithErrorHTTP(w, 406, "User does not exists in the database")
			return
		}
		lib.RespondWithErrorHTTP(w, 500, "Failed to collect user data in the database")
		return
	}
	data := make(map[string]interface{})
	data["firstname"] = user.FirstName
	data["lastname"] = user.LastName
	data["profile_picture"] = user.ProfilePicture
	data["locale"] = user.Locale
	data["available_locales"] = lib.LocaleAvailable
	data["email"] = user.Email
	if user.AccountType.Level == 1 {
		data["username"] = user.Username
	} else {
		data["oauth"] = user.AccountType.Type
	}
	lib.RespondWithJSON(w, http.StatusOK, data)
}
