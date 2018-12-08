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

func insertDataAddComment(filmID string, user1, user2 coltypes.User) {
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
}

func TestAddCommentEmptyContent(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	body, err := lib.InterfaceToByte(map[string]interface{}{
		"content": "",
	})
	if err != nil {
		t.Error(err)
	}
	r := tests.CreateRequest("POST", "/api/v1/comment/"+imdb, body, context)
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
		"error": "Comment content can't be empty",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestAddCommentUnknownMovie(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	body, err := lib.InterfaceToByte(map[string]interface{}{
		"content": "Hello world",
	})
	if err != nil {
		t.Error(err)
	}
	r := tests.CreateRequest("POST", "/api/v1/comment/"+imdb, body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "No movie linked with this filmId",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestAddCommentUnknownUser(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	_ = tests.InsertMovie(coltypes.Movie{IDimdb: imdb}, tests.MongoDB)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	body, err := lib.InterfaceToByte(map[string]interface{}{
		"content": "Hello world",
	})
	if err != nil {
		t.Error(err)
	}
	r := tests.CreateRequest("POST", "/api/v1/comment/"+imdb, body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "UserID does not exists in the database",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestAddCommentInvalidBody(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/api/v1/comment/"+imdb, nil, context)
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

type testFormatJSON struct {
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"createdat"`
	Fullname       string    `json:"fullname"`
	ProfilePicture string    `json:"profile_picture"`
	UserLocale     string    `json:"user_locale"`
}

func TestAddComment(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	user1 := tests.InsertUser(coltypes.User{FirstName: "firstname1", LastName: "1lastname", ProfilePicture: "picture1", Locale: "en"}, tests.MongoDB)
	user2 := tests.InsertUser(coltypes.User{FirstName: "firstname2", LastName: "2lastname", ProfilePicture: "picture2", Locale: "it"}, tests.MongoDB)
	insertDataAddComment(imdb, user1, user2)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  user1.ID,
	}
	body, err := lib.InterfaceToByte(map[string]interface{}{
		"content": "Hello world",
	})
	if err != nil {
		t.Error(err)
	}
	tests.TimeTest = time.Now()
	r := tests.CreateRequest("POST", "/api/v1/comment/"+imdb, body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	if w.Result().StatusCode != 200 {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m.\n", 200, w.Result().StatusCode)
	}
	var response testFormatJSON
	expectedJSONResponse := testFormatJSON{
		Content:        "Hello world",
		CreatedAt:      time.Now(),
		Fullname:       "firstname1 1.",
		ProfilePicture: "picture1",
		UserLocale:     "en",
	}
	err = tests.ChargeResponse(w, &response)
	if err != nil {
		t.Error(err)
	}
	if compare := pretty.Compare(&expectedJSONResponse, response); compare != "" {
		t.Error(compare)
	}
	dbState, err := query.FindComments(bson.M{}, bson.M{"_id": 0}, tests.MongoDB)
	if err != nil {
		t.Errorf("%v", err)
	}
	expectedDatabaseState := []coltypes.Comment{
		coltypes.Comment{
			ID:        "",
			IDimdb:    imdb,
			UserID:    user1.ID,
			Content:   "This is my content 1",
			CreatedAt: time.Date(2017, 1, 2, 10, 1, 0, 0, time.Local),
		},
		coltypes.Comment{
			ID:        "",
			IDimdb:    "filmID2",
			UserID:    user2.ID,
			Content:   "This is my content 2",
			CreatedAt: time.Date(2017, 2, 2, 10, 1, 0, 0, time.Local),
		},
		coltypes.Comment{
			ID:        "",
			IDimdb:    imdb,
			UserID:    user1.ID,
			Content:   "Hello world",
			CreatedAt: time.Now(),
		},
	}
	if compare := pretty.Compare(&expectedDatabaseState, dbState); compare != "" {
		t.Error(compare)
	}
}
