package query

import (
	"errors"

	coltypes "../collections"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func InsertUser(profile coltypes.User, db *mgo.Database) error {
	return db.C("users").Insert(profile)
}

func FindUser(data bson.M, db *mgo.Database) (coltypes.User, error) {
	var user coltypes.User
	if err := db.
		C("users").
		Find(data).
		One(&user); err != nil {
		if err == mgo.ErrNotFound {
			return coltypes.User{}, errors.New("Not Found")
		}
		return coltypes.User{}, err
	}
	return user, nil
}

func FindUsers(data bson.M, db *mgo.Database) ([]coltypes.User, error) {
	var users []coltypes.User
	if err := db.
		C("users").
		Find(data).
		All(&users); err != nil {
		return []coltypes.User{}, err
	}
	return users, nil
}

func UpdateUsers(colQuerier, change bson.M, db *mgo.Database) error {
	err := db.
		C("users").
		Update(colQuerier, bson.M{"$set": change})
	if err != nil {
		return err
	}
	return nil
}
