package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"../../lib"
	"../../mongodb"
	"./auth"
	"./comment"
	"./mail"
	"./movie"
	"./profile"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
)

// HandleAPIRoutes instantiates and populates the router
func handleAPIRoutes() *mux.Router {
	// instantiating the router
	api := mux.NewRouter()

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lib.RespondWithJSON(w, 200, "OK")
	})
	// Don't forget the exception in init.go
	api.HandleFunc("/api/auth/{provider}/callback", auth.OAuthProviderCallback)
	api.HandleFunc("/api/auth/{provider}", auth.OAuthProvider).Methods("GET")
	api.HandleFunc("/api/v1/accounts/register", auth.Register).Methods("POST")
	api.HandleFunc("/api/v1/accounts/login", auth.Login).Methods("POST")
	api.HandleFunc("/api/v1/accounts/resetpassword", auth.ResetPassword).Methods("POST")
	api.HandleFunc("/api/v1/mails/forgotpassword", mail.ForgotPassword).Methods("POST")
	api.HandleFunc("/api/v1/profiles", profile.GetProfile).Methods("GET")
	api.HandleFunc("/api/v1/profiles", profile.EditData).Methods("POST")
	api.HandleFunc("/api/v1/profiles/password", profile.EditPassword).Methods("POST")
	api.HandleFunc("/api/v1/profiles/picture", profile.UploadPicture).Methods("POST")
	api.HandleFunc("/api/v1/users/{username}", profile.GetUser).Methods("GET")
	api.HandleFunc("/api/v1/movies/item/{filmId}/{language}", movie.GetMovie).Methods("GET")
	api.HandleFunc("/api/v1/movies/category/{category}/{offset}/{numberItems}/{language}", movie.GetCategory).Methods("GET")
	api.HandleFunc("/api/v1/movies/search/{offset}/{numberItems}/{language}", movie.SearchMovies).Methods("GET")
	api.HandleFunc("/api/v1/movies/view/{filmId}", movie.AddView).Methods("POST")
	api.HandleFunc("/api/v1/torrents/{filmId}", movie.GetTorrentMagnets).Methods("GET")
	api.HandleFunc("/api/v1/comment/{filmId}", comment.Comments)
	return api
}

func countMovie(db *mgo.Database) int {
	count, _ := db.C("movies").Find(nil).Count()
	return count
}

func main() {
	portPtr := flag.String("port", "3000", "port your want to listen on")
	flag.Parse()
	if *portPtr != "" {
		fmt.Printf("running on port: %s\n", *portPtr)
	}
	var dbName = os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "hypertube"
	}
	var db *mgo.Session
	var connectionError string
	for {
		db, connectionError = mongodb.MongoDBConn(dbName)
		if connectionError == "" {
			break
		}
		fmt.Println(connectionError, ", go retry...")
	}

	count := countMovie(db.DB(dbName))
	if count < 5 {
		fmt.Println(count, " movies, need to launch script")
		output, err := exec.Command("sh", "/app/launchScriptMovies.sh").CombinedOutput()
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}
		fmt.Println(string(output))
	} else {
		fmt.Println(count, " movies, no need to launch script")
	}

	mailjet := lib.MailJetConn()
	auth.InitOAuth2() // Init Omniauth >> Google Plus + API 42
	router := handleAPIRoutes()
	enhancedRouter := enhanceHandlers(router, db, mailjet)
	if err := http.ListenAndServe(":"+*portPtr, enhancedRouter); err != nil {
		log.Fatal(lib.PrettyError(err.Error()))
	}
}
