package movie

import (
	"net/http"
	"os"

	"../../../lib"
	"../../../mongodb/collections"
	"../../../mongodb/query"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type movieVideos struct {
	ID   string `json:"id,omitempty"`
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
	Size int    `json:"size,omitempty"`
}

type movieSimilar struct {
	ID         int    `json:"id,omitempty"`
	Title      string `json:"title,omitempty"`
	PosterPath string `json:"poster_path,omitempty"`
}

type Image struct {
	Path string `json:"file_path,omitempty"`
}

type Images struct {
	Backdrops []Image `json:"backdrops,omitempty"`
}

type trailerResult struct {
	Result []movieVideos `json:"results,omitempty"`
}

type similarResult struct {
	Result []movieSimilar `json:"results,omitempty"`
}

type person struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Character   string `json:"character,omitempty"`
	ProfilePath string `json:"profile_path,omitempty"`
}

type cast struct {
	Cast []person `json:"cast,omitempty"`
}

// Movie is the json data structure for the API response
type Movie struct {
	ID             int           `json:"id,omitempty"`
	IDimdb         string        `json:"imdb_id,omitempty"`
	Title          string        `json:"title,omitempty"`
	Resume         string        `json:"overview,omitempty"`
	Rating         float32       `json:"vote_average,omitempty"`
	PosterPath     string        `json:"poster_path,omitempty"`
	ProductionYear string        `json:"release_date,omitempty"`
	Duration       int           `json:"runtime,omitempty"`
	Trailer        trailerResult `json:"videos,omitempty"`
	Casting        cast          `json:"credits,omitempty"`
	Similar        similarResult `json:"similar,omitempty"`
	Images         Images        `json:"images,omitempty"`
}

func formatOutput(film Movie, filmDB collections.Movie, userId string) map[string]interface{} {
	data := map[string]interface{}{
		"id_imdb":         film.IDimdb,
		"id_tmdb":         film.ID,
		"title":           film.Title,
		"resume":          film.Resume,
		"rating":          filmDB.Rating.Average,
		"votes":           filmDB.Rating.Number,
		"cover":           "https://image.tmdb.org/t/p/w500" + film.PosterPath,
		"production_year": film.ProductionYear,
		"duration":        film.Duration,
		"casting":         []interface{}{},
		"similar":         []interface{}{},
		"images":          []string{},
		"genres":          filmDB.Genres,
		"is_watched":      hasBeenWatched(userId, filmDB.WatchedBy),
		"movie": map[string]interface{}{
			"status":    filmDB.Video.Status,
			"path":      filmDB.Video.Path,
			"stream":    filmDB.Video.Stream,
			"subtitles": filmDB.Video.Subtitles,
		},
	}
	// Trailer
	if len(film.Trailer.Result) >= 1 {
		data["trailer"] = map[string]interface{}{
			"url":     "https://www.youtube.com/watch?v=" + film.Trailer.Result[0].Key,
			"title":   film.Trailer.Result[0].Name,
			"quality": film.Trailer.Result[0].Size,
		}
	}
	// Casting
	for i, person := range film.Casting.Cast {
		if person.ProfilePath != "" {
			data["casting"] = append(data["casting"].([]interface{}),
				map[string]interface{}{
					"name":        person.Name,
					"character":   person.Character,
					"picturePath": "https://image.tmdb.org/t/p/w500" + person.ProfilePath,
				})
		}
		if i >= 15 {
			break
		}
	}
	// Similar movies
	for _, movie := range film.Similar.Result {
		data["similar"] = append(data["similar"].([]interface{}),
			map[string]interface{}{
				"id":        movie.ID,
				"title":     movie.Title,
				"coverPath": "https://image.tmdb.org/t/p/w500" + movie.PosterPath,
			})
	}
	// Images
	var images []string
	for _, image := range film.Images.Backdrops {
		images = append(images, "https://image.tmdb.org/t/p/w500"+image.Path)
	}
	data["images"] = images
	return data
}

func storeAPIResponse(query string, data interface{}) (int, string, error) {
	resp, err := http.Get(query)
	if err != nil {
		return 500, "Call to external api failed", err
	}
	return lib.ReaderJSONToInterface(resp.Body, &data)
}

var (
	apiKey = os.Getenv("TMDB_APIKEY")
)

// GetMovie is the route '/api/v1/movies/item/{filmId}/{language}' with the method GET.
func GetMovie(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.MongoDB).(*mgo.Database)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	userID, _ := r.Context().Value(lib.UserID).(string)
	vars := mux.Vars(r)
	var (
		filmID   = vars["filmId"] // Can be imdb or tmdb id
		language = vars["language"]
	)
	var film Movie
	errCode, errStatus, err := storeAPIResponse(
		"https://api.themoviedb.org/3/movie/"+filmID+
			"?api_key="+apiKey+
			"&language="+language+
			"&append_to_response=images,credits,videos,similar",
		&film,
	)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errStatus)
		return
	}
	// If no movie resume (no standard language), get the english resume
	if film.Resume == "" {
		var filmOnly Movie
		errCode, errStatus, err := storeAPIResponse(
			"https://api.themoviedb.org/3/movie/"+filmID+"?api_key="+apiKey,
			&filmOnly,
		)
		if err != nil {
			lib.RespondWithErrorHTTP(w, errCode, errStatus)
			return
		}
		film.Resume = filmOnly.Resume
	}
	// Get database data -> Rating/Comments
	dbMovie, err := query.FindMovie(bson.M{"_id": film.IDimdb}, db)
	if err != nil {
		if err.Error() == "Not Found" {
			lib.RespondWithErrorHTTP(w, 406, "Movie does not exists in the database")
			return
		}
		lib.RespondWithErrorHTTP(w, 500, "Failed to collect movie data in the database")
		return
	}
	lib.RespondWithJSON(w, http.StatusOK, formatOutput(film, dbMovie, userID))
}
