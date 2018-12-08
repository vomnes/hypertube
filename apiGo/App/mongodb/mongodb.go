package mongodb

import (
	"log"
	"os"

	"../lib"
	mgo "gopkg.in/mgo.v2"
)

const (
	// host = "localhost"
	// host = "0.0.0.0"
	port = "27017"
	// port = "1234"
)

// MongoDBConn allows to create a connection with MongoDB database
func MongoDBConn(dbName string) (*mgo.Session, string) {
	var host = os.Getenv("MONGO_DB_HOST")
	if host == "" {
		host = "localhost"
	}
	if dbName == "" {
		dbName = "db_hypertube_tests"
	}
	dbURL := "mongodb://" + host + ":" + port + "/" + dbName
	session, err := mgo.Dial(dbURL)
	if err != nil {
		if err.Error() == "no reachable servers" {
			return session, "Cannot connect to mongo: " + dbURL + " - " + err.Error()
		} else {
			log.Fatal(lib.PrettyError("Cannot connect to mongo: " + dbURL + " - " + err.Error()))
		}
	}
	return session, ""
}
