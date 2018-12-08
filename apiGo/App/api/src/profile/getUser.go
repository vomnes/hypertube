package profile

import (
	"net/http"

	"../../../lib"
	coltype "../../../mongodb/collections"
	"../../../mongodb/query"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.MongoDB).(*mgo.Database)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	vars := mux.Vars(r)
	targetUsername := vars["username"]
	if !lib.IsValidUsername(targetUsername) {
		lib.RespondWithErrorHTTP(w, 406, "Target username is invalid")
		return
	}
	// Get user data
	var user coltype.User
	user, err := query.FindUser(bson.M{"username": targetUsername}, db)
	if err != nil {
		if err.Error() == "Not Found" {
			lib.RespondWithErrorHTTP(w, 406, "User does not exists in the database")
			return
		}
		lib.RespondWithErrorHTTP(w, 500, "Failed to collect user data in the database")
		return
	}
	lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"username":        user.Username,
		"firstname":       user.FirstName,
		"lastname":        user.LastName,
		"profile_picture": user.ProfilePicture,
	})
}
