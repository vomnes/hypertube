package comment

import (
	"log"
	"net/http"

	"../../../lib"
	"../../../mongodb/query"
	"gopkg.in/mgo.v2/bson"
)

type body struct {
	ID string `json:"comment_id"`
}

func (d *data) deleteComment(w http.ResponseWriter, r *http.Request) {
	var inputData body
	errCode, errContent, err := lib.ReaderJSONToInterface(r.Body, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	if inputData.ID == "" {
		lib.RespondWithErrorHTTP(w, 400, "CommentID must be defined")
		return
	}
	err = query.DeleteComment(
		bson.M{
			"_id":    inputData.ID,
			"idimdb": d.FilmID,
			"userid": d.UserID,
		},
		d.DB)
	if err != nil {
		if err.Error() == "not found" {
			// Either the comment is not found, not from the connected user or
			// the filmId is not linked with the comment id
			lib.RespondWithErrorHTTP(w, 401, "Unauthorized to remove this comment")
			return
		}
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to remove comment " + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Remove data failed")
		return
	}
	lib.RespondEmptyHTTP(w, http.StatusOK)
}
