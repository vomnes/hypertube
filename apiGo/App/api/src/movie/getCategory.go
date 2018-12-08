package movie

import (
	"net/http"
	"strconv"
	"strings"

	"../../../lib"
	coltypes "../../../mongodb/collections"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GetCategory is the route '/api/v1/movies/category/{category}/{offset}/{numberItems}/{language}' with the method GET.
func GetCategory(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.MongoDB).(*mgo.Database)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	userID, _ := r.Context().Value(lib.UserID).(string)
	vars := mux.Vars(r)
	language := vars["language"]
	offset, _ := strconv.Atoi(vars["offset"])
	numberItems, err := strconv.Atoi(vars["numberItems"])
	if err != nil {
		numberItems = 20
	}
	request := bson.M{"genres": bson.M{"$in": []string{strings.Title(vars["category"])}}}
	if vars["category"] == "popular" {
		request = bson.M{"rating.average": bson.M{"$gte": 7.5}}
	}
	var movies []coltypes.Movie
	err = db.
		C("movies").
		Find(request).
		Select(bson.M{
			"_id":            1,
			"poster":         1,
			"title":          1,
			"language_title": 1,
			"rating.average": 1,
			"watchedby":      1,
			"video.status":   1,
		}).
		Sort("-year", "-rating.average").
		Skip(offset).
		Limit(numberItems).
		All(&movies)
	if err != nil {
		lib.RespondWithErrorHTTP(w, 500, "Failed to collect movies data in the database")
		return
	}
	f := formatMovie{
		movies:   movies,
		userID:   userID,
		language: language,
	}
	response := f.formatListMovies()
	if len(response) > 0 {
		lib.RespondWithJSON(w, http.StatusOK, response)
		return
	}
	lib.RespondWithJSON(w, http.StatusOK, map[string]string{
		"status": "No (more) data",
	})
}
