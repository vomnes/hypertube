package comment

import (
	"net/http"

	"../../../lib"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (d *data) getComments(w http.ResponseWriter, r *http.Request) {
	documentPipeline := []bson.M{
		// Select only the comments with the right filmID
		bson.M{
			"$match": bson.M{
				"idimdb": d.FilmID,
			},
		},
		// Keep only the elements interesting for our request
		bson.M{
			"$project": bson.M{
				"_id":       1,
				"userid":    1,
				"content":   1,
				"createdat": 1,
			},
		},
		// Collect the users data liked with film.userid
		bson.M{
			"$lookup": bson.M{
				"from":         "users",
				"foreignField": "_id",
				"localField":   "userid",
				"as":           "user",
			},
		},
		// Transform user from object array to object
		bson.M{
			"$unwind": "$user",
		},
		// Transform the object to get almost our final form
		bson.M{
			"$project": bson.M{
				"_id": 0,
				"id": bson.M{
					// Shows the comment id only if the connected user is the owner
					"$cond": bson.M{
						"if": bson.M{
							"$eq": []interface{}{
								"$userid",
								d.UserID,
							},
						},
						"then": "$_id",
						"else": "$REMOVE",
					},
				},
				"content":   1,
				"createdat": 1,
				"fullname": bson.M{
					"$concat": []interface{}{
						"$user.firstname",
						" ",
						bson.M{
							"$substr": []interface{}{
								"$user.lastname",
								0,
								1,
							},
						},
						".",
					},
				},
				"profile_picture": "$user.profile_picture",
				"user_locale":     "$user.locale",
			},
		},
	}
	var comments []map[string]interface{}
	if err := d.DB.C("comments").Pipe(documentPipeline).All(&comments); err != nil {
		if err == mgo.ErrNotFound {
			lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{})
		}
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	if len(comments) == 0 {
		lib.RespondWithJSON(w, http.StatusOK, map[string]string{
			"status": "No comments",
		})
		return
	}
	lib.RespondWithJSON(w, http.StatusOK, comments)
}
