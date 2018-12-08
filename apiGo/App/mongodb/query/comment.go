package query

import (
	"errors"

	coltypes "../collections"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func InsertComment(comment coltypes.Comment, db *mgo.Database) error {
	return db.C("comments").Insert(comment)
}

func FindComment(data bson.M, db *mgo.Database) (coltypes.Comment, error) {
	var comment coltypes.Comment
	if err := db.
		C("comments").
		Find(data).
		One(&comment); err != nil {
		if err == mgo.ErrNotFound {
			return coltypes.Comment{}, errors.New("Not Found")
		}
		return coltypes.Comment{}, err
	}
	return comment, nil
}

func FindComments(data, selectData bson.M, db *mgo.Database) ([]coltypes.Comment, error) {
	var comments []coltypes.Comment
	if err := db.
		C("comments").
		Find(data).
		Select(selectData).
		All(&comments); err != nil {
		return []coltypes.Comment{}, err
	}
	return comments, nil
}

func DeleteComment(identifier bson.M, db *mgo.Database) error {
	return db.C("comments").Remove(identifier)
}
