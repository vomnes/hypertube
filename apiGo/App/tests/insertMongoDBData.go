package tests

import (
	"log"

	"../lib"
	coltype "../mongodb/collections"
	mgoQuery "../mongodb/query"
	mgo "gopkg.in/mgo.v2"
)

func InsertUser(user coltype.User, db *mgo.Database) coltype.User {
	// Generate ID
	if user.ID == "" {
		user.ID = lib.GetRandomString(42)
	}
	err := mgoQuery.InsertUser(user, db)
	if err != nil {
		log.Fatal(lib.PrettyError("[TEST] Failed to insert user data" + err.Error()))
	}
	return user
}

func InsertMovie(movie coltype.Movie, db *mgo.Database) coltype.Movie {
	// Generate ID
	if movie.IDimdb == "" {
		movie.IDimdb = "tttest" + lib.GetRandomString(42)
	}
	err := mgoQuery.InsertMovie(movie, db)
	if err != nil {
		log.Fatal(lib.PrettyError("[TEST] Failed to insert movie data" + err.Error()))
	}
	return movie
}

func InsertComment(comment coltype.Comment, db *mgo.Database) coltype.Comment {
	// Generate ID
	if comment.ID == "" {
		comment.ID = "test" + lib.GetRandomString(42)
	}
	err := mgoQuery.InsertComment(comment, db)
	if err != nil {
		log.Fatal(lib.PrettyError("[TEST] Failed to insert comment data" + err.Error()))
	}
	return comment
}
