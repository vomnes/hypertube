package comment

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"../../../lib"
	coltypes "../../../mongodb/collections"
	"../../../mongodb/query"
	"../../../tests"
	"github.com/kylelemons/godebug/pretty"
	"gopkg.in/mgo.v2/bson"
)

func insertDataDeleteComment(filmID string, user1, user2 coltypes.User) {
	_ = tests.InsertMovie(coltypes.Movie{IDimdb: filmID}, tests.MongoDB)
	_ = tests.InsertComment(coltypes.Comment{
		ID:        "1",
		IDimdb:    filmID,
		UserID:    user1.ID,
		Content:   "This is my content 1",
		CreatedAt: time.Date(2017, 1, 2, 10, 1, 0, 0, time.Local),
	}, tests.MongoDB)
	_ = tests.InsertComment(coltypes.Comment{
		ID:        "2",
		IDimdb:    "filmID2",
		UserID:    user2.ID,
		Content:   "This is my content 2",
		CreatedAt: time.Date(2017, 2, 2, 10, 1, 0, 0, time.Local),
	}, tests.MongoDB)
	_ = tests.InsertComment(coltypes.Comment{
		ID:        "3",
		IDimdb:    filmID,
		UserID:    user1.ID,
		Content:   "This is my content 3",
		CreatedAt: time.Date(2018, 11, 15, 9, 13, 0, 0, time.Local),
	}, tests.MongoDB)
}

func TestDeleteCommentEmptyID(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	body, err := lib.InterfaceToByte(map[string]interface{}{
		"comment_id": "",
	})
	if err != nil {
		t.Error(err)
	}
	r := tests.CreateRequest("DELETE", "/api/v1/comment/"+imdb, body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 400, map[string]interface{}{
		"error": "CommentID must be defined",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestDeleteCommentInvalidBody(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("DELETE", "/api/v1/comment/"+imdb, nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if !strings.Contains(output, "Failed to decode json readerEOF") {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Failed to decode JSON reader",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestDeleteComment(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	user1 := tests.InsertUser(coltypes.User{FirstName: "firstname1", LastName: "1lastname", ProfilePicture: "picture1", Locale: "en"}, tests.MongoDB)
	user2 := tests.InsertUser(coltypes.User{FirstName: "firstname2", LastName: "2lastname", ProfilePicture: "picture2", Locale: "it"}, tests.MongoDB)
	insertDataDeleteComment(imdb, user1, user2)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  user1.ID,
	}
	body, err := lib.InterfaceToByte(map[string]interface{}{
		"comment_id": "3",
	})
	if err != nil {
		t.Error(err)
	}
	dbState, err := query.FindComments(bson.M{}, nil, tests.MongoDB)
	if err != nil {
		t.Errorf("%v", err)
	}
	expectedDatabaseState := []coltypes.Comment{
		coltypes.Comment{
			ID:        "1",
			IDimdb:    imdb,
			UserID:    user1.ID,
			Content:   "This is my content 1",
			CreatedAt: time.Date(2017, 1, 2, 10, 1, 0, 0, time.Local),
		},
		coltypes.Comment{
			ID:        "2",
			IDimdb:    "filmID2",
			UserID:    user2.ID,
			Content:   "This is my content 2",
			CreatedAt: time.Date(2017, 2, 2, 10, 1, 0, 0, time.Local),
		},
		coltypes.Comment{
			ID:        "3",
			IDimdb:    imdb,
			UserID:    user1.ID,
			Content:   "This is my content 3",
			CreatedAt: time.Date(2018, 11, 15, 9, 13, 0, 0, time.Local),
		},
	}
	if compare := pretty.Compare(&expectedDatabaseState, dbState); compare != "" {
		t.Error("Before delete:", compare)
	}
	r := tests.CreateRequest("DELETE", "/api/v1/comment/"+imdb, body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	dbState, err = query.FindComments(bson.M{}, nil, tests.MongoDB)
	if err != nil {
		t.Errorf("%v", err)
	}
	expectedDatabaseState = []coltypes.Comment{
		coltypes.Comment{
			ID:        "1",
			IDimdb:    imdb,
			UserID:    user1.ID,
			Content:   "This is my content 1",
			CreatedAt: time.Date(2017, 1, 2, 10, 1, 0, 0, time.Local),
		},
		coltypes.Comment{
			ID:        "2",
			IDimdb:    "filmID2",
			UserID:    user2.ID,
			Content:   "This is my content 2",
			CreatedAt: time.Date(2017, 2, 2, 10, 1, 0, 0, time.Local),
		},
	}
	if compare := pretty.Compare(&expectedDatabaseState, dbState); compare != "" {
		t.Error("After delete:", compare)
	}
}

func TestDeleteCommentUnauthorized(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	user1 := tests.InsertUser(coltypes.User{FirstName: "firstname1", LastName: "1lastname", ProfilePicture: "picture1", Locale: "en"}, tests.MongoDB)
	user2 := tests.InsertUser(coltypes.User{FirstName: "firstname2", LastName: "2lastname", ProfilePicture: "picture2", Locale: "it"}, tests.MongoDB)
	insertDataDeleteComment(imdb, user1, user2)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  user1.ID,
	}
	body, err := lib.InterfaceToByte(map[string]interface{}{
		"comment_id": "2",
	})
	if err != nil {
		t.Error(err)
	}
	dbState, err := query.FindComments(bson.M{}, nil, tests.MongoDB)
	if err != nil {
		t.Errorf("%v", err)
	}
	expectedDatabaseState := []coltypes.Comment{
		coltypes.Comment{
			ID:        "1",
			IDimdb:    imdb,
			UserID:    user1.ID,
			Content:   "This is my content 1",
			CreatedAt: time.Date(2017, 1, 2, 10, 1, 0, 0, time.Local),
		},
		coltypes.Comment{
			ID:        "2",
			IDimdb:    "filmID2",
			UserID:    user2.ID,
			Content:   "This is my content 2",
			CreatedAt: time.Date(2017, 2, 2, 10, 1, 0, 0, time.Local),
		},
		coltypes.Comment{
			ID:        "3",
			IDimdb:    imdb,
			UserID:    user1.ID,
			Content:   "This is my content 3",
			CreatedAt: time.Date(2018, 11, 15, 9, 13, 0, 0, time.Local),
		},
	}
	if compare := pretty.Compare(&expectedDatabaseState, dbState); compare != "" {
		t.Error("Before delete:", compare)
	}
	r := tests.CreateRequest("DELETE", "/api/v1/comment/"+imdb, body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 401, map[string]interface{}{
		"error": "Unauthorized to remove this comment",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	dbState, err = query.FindComments(bson.M{}, nil, tests.MongoDB)
	if err != nil {
		t.Errorf("%v", err)
	}
	expectedDatabaseState = []coltypes.Comment{
		coltypes.Comment{
			ID:        "1",
			IDimdb:    imdb,
			UserID:    user1.ID,
			Content:   "This is my content 1",
			CreatedAt: time.Date(2017, 1, 2, 10, 1, 0, 0, time.Local),
		},
		coltypes.Comment{
			ID:        "2",
			IDimdb:    "filmID2",
			UserID:    user2.ID,
			Content:   "This is my content 2",
			CreatedAt: time.Date(2017, 2, 2, 10, 1, 0, 0, time.Local),
		},
		coltypes.Comment{
			ID:        "3",
			IDimdb:    imdb,
			UserID:    user1.ID,
			Content:   "This is my content 3",
			CreatedAt: time.Date(2018, 11, 15, 9, 13, 0, 0, time.Local),
		},
	}
	if compare := pretty.Compare(&expectedDatabaseState, dbState); compare != "" {
		t.Error("After delete:", compare)
	}
}
