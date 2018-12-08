package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"../../lib"
	"github.com/gorilla/mux"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
	"github.com/markbates/goth"
	"github.com/rs/cors"
	mgo "gopkg.in/mgo.v2"
)

type adapter func(http.Handler) http.Handler

// adapt transforms an handler without changing it's type. Usefull for authentification.
func adapt(h http.Handler, adapters ...adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// adapt the request by checking the auth and filling the context with usefull data
func enhanceHandlers(r *mux.Router, db *mgo.Session, mailjet *mailjet.Client) http.Handler {
	return adapt(r, withRights(), withConnections(db, mailjet))
	// return adapt(r, withRights(), withConnections(db, mailjet), withCors())
}

// withConnections is an adapter that copy the access to the database to serve a specific call
func withConnections(db *mgo.Session, mailjet *mailjet.Client) adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dbsession := db.Copy()
			defer dbsession.Close() // cleaning up
			db_name := os.Getenv("MONGO_DB_NAME")
			if db_name == "" {
				db_name = "db_hypertube_tests"
			}
			ctx := context.WithValue(r.Context(), lib.MongoDB, dbsession.DB(db_name))
			ctx = context.WithValue(ctx, lib.MongoDBSession, dbsession)
			ctx = context.WithValue(ctx, lib.MailJet, mailjet)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// withRights is an adapter that verify the user exists, verify the token,
// and attach userId and username to the request.
func withRights() adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			routeURL := *r.URL
			routeURLString := routeURL.String()
			noCheckJWT := []string{
				"/",
				"/api/v1/accounts/register",
				"/api/v1/accounts/login",
				"/api/v1/accounts/resetpassword",
				"/api/v1/mails/forgotpassword",
			}
			if lib.StringInArray(routeURLString, noCheckJWT) {
				h.ServeHTTP(w, r)
				return
			}
			for oAuthProvider, _ := range goth.GetProviders() {
				if routeURLString == "/api/auth/"+oAuthProvider ||
					strings.Contains(routeURLString, "/api/auth/"+oAuthProvider+"/callback") {
					h.ServeHTTP(w, r)
					return
				}
			}
			var tokenString string
			// Get token from the Authorization header format: Authorization: Bearer <jwt>
			tokens := r.Header.Get("Authorization")
			if tokens != "" {
				tokenString = tokens
				if !strings.HasPrefix(tokenString, "Bearer ") {
					lib.RespondWithErrorHTTP(w, 403, "Access denied - Authorization wrong standard")
					return
				}
				tokenString = strings.TrimPrefix(tokenString, "Bearer ")
			} else {
				lib.RespondWithErrorHTTP(w, 403, "Access denied")
				return
			}
			// Check JWT validity on every request
			// Parse takes the token string and a function for looking up the key
			claims, err := lib.AnalyseJWT(tokenString)
			if err != nil {
				lib.RespondWithErrorHTTP(w, 403, "Access denied - "+err.Error())
				return
			}
			if claims["username"] == nil || claims["userId"] == nil ||
				claims["firstname"] == nil || claims["lastname"] == nil ||
				claims["profile_picture"] == nil {
				lib.RespondWithErrorHTTP(w, 403, "Access denied - Not the right data in JWT")
				return
			}
			// Attach data from the token to the request
			ctx := context.WithValue(r.Context(), lib.UserID, claims["userId"])
			ctx = context.WithValue(ctx, lib.Username, claims["username"])
			ctx = context.WithValue(ctx, lib.FirstName, claims["firstname"])
			ctx = context.WithValue(ctx, lib.LastName, claims["lastname"])
			ctx = context.WithValue(ctx, lib.ProfilePicture, claims["profile_picture"])
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// withCors is an adpater that allowed the specific headers we need for our requests from a
// different domain.
func withCors() adapter {
	return func(h http.Handler) http.Handler {
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost"},
			AllowedHeaders:   []string{""},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
			AllowCredentials: true,
		})
		c = cors.AllowAll()
		return c.Handler(h)
	}
}
