package query

import (
	"errors"

	coltypes "../collections"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func InsertMovie(movie coltypes.Movie, db *mgo.Database) error {
	return db.C("movies").Insert(movie)
}

func FindMovie(data bson.M, db *mgo.Database) (coltypes.Movie, error) {
	var movie coltypes.Movie
	if err := db.
		C("movies").
		Find(data).
		One(&movie); err != nil {
		if err == mgo.ErrNotFound {
			return coltypes.Movie{}, errors.New("Not Found")
		}
		return coltypes.Movie{}, err
	}
	return movie, nil
}

func FindMovies(data, selectData bson.M, db *mgo.Database) ([]coltypes.Movie, error) {
	var movies []coltypes.Movie
	if err := db.
		C("movies").
		Find(data).
		Select(selectData).
		All(&movies); err != nil {
		return []coltypes.Movie{}, err
	}
	return movies, nil
}

func UpdateMovie(colQuerier, change bson.M, db *mgo.Database) error {
	err := db.
		C("movies").
		Update(colQuerier, bson.M{"$set": change})
	if err != nil {
		return err
	}
	return nil
}
