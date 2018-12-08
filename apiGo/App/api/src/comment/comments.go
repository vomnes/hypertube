package comment

import (
	"net/http"

	"../../../lib"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
)

type data struct {
	DB     *mgo.Database
	UserID string
	FilmID string
}

// Comments allows to get, add and remove the comment(s)
func Comments(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.MongoDB).(*mgo.Database)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	if !lib.CheckHTTPMethod(r.Method, []string{"GET", "POST", "DELETE"}) {
		lib.RespondWithErrorHTTP(w, 404, "Page not found")
		return
	}
	userID, _ := r.Context().Value(lib.UserID).(string)
	vars := mux.Vars(r)
	filmID := vars["filmId"]
	d := data{db, userID, filmID}
	switch r.Method {
	case "GET":
		d.getComments(w, r)
		return
	case "POST":
		d.addComment(w, r)
		return
	case "DELETE":
		d.deleteComment(w, r)
		return
	}
}

func setFullname(firstname, lastname string) string {
	if firstname != "" && lastname != "" {
		return firstname + " " + string(lastname[0]) + "."
	}
	return ""
}
