package movie

import (
	"net/http"
	"strconv"

	"../../../lib"
	coltypes "../../../mongodb/collections"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type limits struct {
	Max float32 `json:"max"`
	Min float32 `json:"min"`
}

type search struct {
	Search string `json:"search"`
	Year   limits `json:"year"`
	Rating limits `json:"rating"`
	Status string `json:"status"`
}

// SearchMovies is the route '/api/v1/movies/search/{offset}/{numberItems}/{language}' with the method GET.
func SearchMovies(w http.ResponseWriter, r *http.Request) {
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
	strSearchParameters := "e30="
	searchParameters, right := r.Header["Search-Parameters"]
	if right && len(searchParameters) > 0 {
		strSearchParameters = searchParameters[0]
	}
	var inputData search
	err = lib.ExtractBase64Struct(strSearchParameters, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, http.StatusNotAcceptable, "Failed to extract base64 search parameters in header")
		return
	}
	var movies []coltypes.Movie
	regexQuery := bson.M{"$regex": ".*" + inputData.Search + ".*", "$options": "i"}
	request := bson.M{
		"$or": []bson.M{
			bson.M{"title": regexQuery},
			bson.M{"language_title." + language: regexQuery},
		},
		"year": bson.M{
			"$gte": inputData.Year.Min,
			"$lte": inputData.Year.Max,
		},
		"rating.average": bson.M{
			"$gte": inputData.Rating.Min,
			"$lte": inputData.Rating.Max,
		},
	}
	if lib.StringInArray(inputData.Status, []string{"not downloaded", "waiting", "downloading", "ready"}) {
		request["video.status"] = bson.M{
			"$in": encodeVideoStatus(inputData.Status),
		}
	}
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
			"genres":         1,
			"year":           1,
			"duration":       1,
		}).
		Sort("-rating.average", "-year"). // Sort by rating (higher first) and then by year
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
		show: show{
			genres:   true,
			year:     true,
			duration: true,
		},
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
