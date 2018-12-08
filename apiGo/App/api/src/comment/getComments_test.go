package comment

import (
	"net/http/httptest"
	"testing"
	"time"

	coltypes "../../../mongodb/collections"
	"../../../tests"
)

func insertData(filmID string) string {
	user1 := tests.InsertUser(coltypes.User{FirstName: "firstname1", LastName: "1lastname", ProfilePicture: "picture1", Locale: "en"}, tests.MongoDB)
	user2 := tests.InsertUser(coltypes.User{FirstName: "firstname2", LastName: "2lastname", ProfilePicture: "picture2", Locale: "it"}, tests.MongoDB)
	user3 := tests.InsertUser(coltypes.User{FirstName: "firstname3", ProfilePicture: "picture3", Locale: "fr"}, tests.MongoDB)
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
		UserID:    user3.ID,
		Content:   "This is my content 3",
		CreatedAt: time.Date(2017, 3, 2, 10, 1, 0, 0, time.Local),
	}, tests.MongoDB)
	_ = tests.InsertComment(coltypes.Comment{
		ID:        "4",
		IDimdb:    filmID,
		UserID:    user2.ID,
		Content:   "This is my content 4",
		CreatedAt: time.Date(2017, 4, 2, 10, 1, 0, 0, time.Local),
	}, tests.MongoDB)
	_ = tests.InsertComment(coltypes.Comment{
		ID:        "5",
		IDimdb:    "filmID4",
		UserID:    user3.ID,
		Content:   "This is my content 5",
		CreatedAt: time.Date(2017, 5, 2, 10, 1, 0, 0, time.Local),
	}, tests.MongoDB)
	_ = tests.InsertComment(coltypes.Comment{
		ID:        "6",
		IDimdb:    filmID,
		UserID:    user1.ID,
		Content:   "This is my content 6",
		CreatedAt: time.Date(2017, 6, 2, 10, 1, 0, 0, time.Local),
	}, tests.MongoDB)
	return user1.ID
}

func TestGetComments(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	userId := insertData(imdb)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  userId,
	}
	r := tests.CreateRequest("GET", "/api/v1/comment/"+imdb, nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, []map[string]interface{}{
		map[string]interface{}{
			"id":              "1",
			"content":         "This is my content 1",
			"createdat":       time.Date(2017, 1, 2, 10, 1, 0, 0, time.Local).Format(time.RFC3339),
			"fullname":        "firstname1 1.",
			"profile_picture": "picture1",
			"user_locale":     "en",
		},
		map[string]interface{}{
			"content":         "This is my content 3",
			"createdat":       time.Date(2017, 3, 2, 10, 1, 0, 0, time.Local).Format(time.RFC3339),
			"fullname":        "firstname3 .",
			"profile_picture": "picture3",
			"user_locale":     "fr",
		},
		map[string]interface{}{
			"content":         "This is my content 4",
			"createdat":       time.Date(2017, 4, 2, 10, 1, 0, 0, time.Local).Format(time.RFC3339),
			"fullname":        "firstname2 2.",
			"profile_picture": "picture2",
			"user_locale":     "it",
		},
		map[string]interface{}{
			"id":              "6",
			"content":         "This is my content 6",
			"createdat":       time.Date(2017, 6, 2, 10, 1, 0, 0, time.Local).Format(time.RFC3339),
			"fullname":        "firstname1 1.",
			"profile_picture": "picture1",
			"user_locale":     "en",
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetCommentsUnkonwnIMDBID(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	insertData(imdb)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("GET", "/api/v1/comment/"+"unknown", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{
		"status": "No comments",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetCommentsUnkonwn(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	insertData(imdb)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("GET", "/api/v1/comment/"+"unknown", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{
		"status": "No comments",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetCommentsOne(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	insertData(imdb)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("GET", "/api/v1/comment/"+"filmID4", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, []map[string]interface{}{
		map[string]interface{}{
			"content":         "This is my content 5",
			"createdat":       time.Date(2017, 5, 2, 10, 1, 0, 0, time.Local).Format(time.RFC3339),
			"fullname":        "firstname3 .",
			"profile_picture": "picture3",
			"user_locale":     "fr",
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}
