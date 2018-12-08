package movie

import (
	"log"
	"net/http"
	"time"

	"../../../lib"
	coltypes "../../../mongodb/collections"
	"../../../mongodb/query"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AddView is the route '/api/v1/movies/view/{filmId}' with the method POST.
func AddView(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.MongoDB).(*mgo.Database)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	userID, _ := r.Context().Value(lib.UserID).(string)
	vars := mux.Vars(r)
	filmID := vars["filmId"]
	movie, err := query.FindMovie(bson.M{"_id": filmID}, db)
	if err != nil {
		if err.Error() == "Not Found" {
			lib.RespondWithErrorHTTP(w, 406, "No movie linked with this filmId")
			return
		}
		log.Println(lib.PrettyError("[DB REQUEST - FIND] Failed to check movie data" + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Failed to check if the film id exists")
		return
	}
	_, err = query.FindUser(bson.M{"_id": userID}, db)
	if err != nil {
		if err.Error() == "Not Found" {
			lib.RespondWithErrorHTTP(w, 406, "User does not exists in the database")
			return
		}
		log.Println(lib.PrettyError("[DB REQUEST - FIND] Failed to check user data" + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Failed to get user data")
		return
	}
	// Update watchedBy
	movie.WatchedBy = append(movie.WatchedBy, coltypes.Watched{
		UserID:    userID,
		WatchedAt: time.Now(),
	})
	err = query.UpdateMovie(bson.M{"_id": filmID}, bson.M{"watchedby": movie.WatchedBy}, db)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - UPDATE] Failed to insert 'movie watched' in the database" + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Failed to update data")
		return
	}
	lib.RespondEmptyHTTP(w, http.StatusOK)
}
