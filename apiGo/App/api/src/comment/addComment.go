package comment

import (
	"log"
	"net/http"
	"time"

	"../../../lib"
	coltypes "../../../mongodb/collections"
	"../../../mongodb/query"
	"gopkg.in/mgo.v2/bson"
)

type comment struct {
	Content string `json:"content"`
}

func (d *data) addComment(w http.ResponseWriter, r *http.Request) {
	var inputData comment
	errCode, errContent, err := lib.ReaderJSONToInterface(r.Body, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	if inputData.Content == "" {
		lib.RespondWithErrorHTTP(w, 400, "Comment content can't be empty")
		return
	}
	_, err = query.FindMovie(bson.M{"_id": d.FilmID}, d.DB)
	if err != nil {
		if err.Error() == "Not Found" {
			lib.RespondWithErrorHTTP(w, 406, "No movie linked with this filmId")
			return
		}
		log.Println(lib.PrettyError("[DB REQUEST - FIND] Failed to check movie data" + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Failed to check if the film id exists")
		return
	}
	user, err := query.FindUser(bson.M{"_id": d.UserID}, d.DB)
	if err != nil {
		if err.Error() == "Not Found" {
			lib.RespondWithErrorHTTP(w, 406, "UserID does not exists in the database")
			return
		}
		log.Println(lib.PrettyError("[DB REQUEST - FIND] Failed to get user data" + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Failed to get user data")
		return
	}
	comment := coltypes.Comment{
		ID:        lib.GetRandomString(42),
		IDimdb:    d.FilmID,
		UserID:    d.UserID,
		Content:   inputData.Content,
		CreatedAt: time.Now(),
	}
	err = query.InsertComment(comment, d.DB)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to insert comment" + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Insert data failed")
		return
	}
	lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"id":              comment.ID,
		"content":         comment.Content,
		"createdat":       comment.CreatedAt,
		"fullname":        setFullname(user.FirstName, user.LastName),
		"profile_picture": user.ProfilePicture,
		"user_locale":     user.Locale,
	})
}
