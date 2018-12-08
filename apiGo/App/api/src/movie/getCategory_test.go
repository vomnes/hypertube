package movie

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"../../../lib"
	coltypes "../../../mongodb/collections"
	"../../../tests"
)

func insertData(category, userID string) {
	_ = tests.InsertMovie(coltypes.Movie{
		IDimdb:        "tt1",
		OriginalTitle: "This is 1",
		Titles: map[string]string{
			"it": "è 1",
			"fr": "C'est le 1",
		},
		Year:   2015,
		Genres: []string{category, "a", "b", "c"},
		Rating: coltypes.Rating{
			Average: 7.5,
		},
		Poster: "Poster 1",
		Video: coltypes.Video{
			Status: -1,
		},
		WatchedBy: []coltypes.Watched{
			coltypes.Watched{
				UserID:    "test_" + lib.GetRandomString(43),
				WatchedAt: time.Date(2018, 2, 2, 10, 1, 0, 0, time.UTC),
			},
		},
	}, tests.MongoDB)
	_ = tests.InsertMovie(coltypes.Movie{
		IDimdb:        "tt9",
		OriginalTitle: "This is 9",
		Year:          2015,
		Genres:        []string{"x", category, "y"},
		Rating: coltypes.Rating{
			Average: 2.5,
		},
		Poster: "Poster 9",
		Video: coltypes.Video{
			Status: 2,
		},
		WatchedBy: []coltypes.Watched{
			coltypes.Watched{
				UserID:    userID,
				WatchedAt: time.Date(2018, 2, 2, 10, 1, 0, 0, time.UTC),
			},
		},
	}, tests.MongoDB)
	_ = tests.InsertMovie(coltypes.Movie{
		IDimdb:        "tt2",
		OriginalTitle: "This is 2",
		Titles: map[string]string{
			"fr": "C'est le 2",
		},
		Year:   2015,
		Genres: []string{"x", "y", category, "c"},
		Rating: coltypes.Rating{
			Average: 9,
		},
		Poster: "Poster 2",
		Video: coltypes.Video{
			Status: 2,
		},
		WatchedBy: []coltypes.Watched{
			coltypes.Watched{
				UserID:    "test_" + lib.GetRandomString(43),
				WatchedAt: time.Date(2017, 2, 3, 10, 1, 0, 0, time.UTC),
			},
			coltypes.Watched{
				UserID:    userID,
				WatchedAt: time.Date(2018, 2, 3, 10, 1, 0, 0, time.UTC),
			},
			coltypes.Watched{
				UserID:    "test_" + lib.GetRandomString(43),
				WatchedAt: time.Date(2018, 2, 5, 10, 1, 0, 0, time.UTC),
			},
		},
	}, tests.MongoDB)
	_ = tests.InsertMovie(coltypes.Movie{
		IDimdb:        "tt3",
		OriginalTitle: "This is 3",
		Year:          2016,
		Genres:        []string{category},
		Rating: coltypes.Rating{
			Average: 1,
		},
		Poster: "Poster 3",
		Video: coltypes.Video{
			Status: 1,
		},
	}, tests.MongoDB)
	_ = tests.InsertMovie(coltypes.Movie{
		IDimdb:        "tt4",
		OriginalTitle: "This is 4",
		Year:          2015,
		Genres:        []string{"a"},
		Rating: coltypes.Rating{
			Average: 1,
		},
		Poster: "Poster 4",
		Video: coltypes.Video{
			Status: 0,
		},
	}, tests.MongoDB)
}

func TestGetCategory(t *testing.T) {
	tests.DbClean()
	category := "Testcategory"
	userID := "test_" + lib.GetRandomString(43)
	insertData(category, userID)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  userID,
	}
	r := tests.CreateRequest("GET", "/api/v1/movies/category/"+strings.ToLower(category)+"/0/20/fr", []byte(`{}`), context)
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
			"id":         "tt3",
			"rating":     1,
			"title":      "This is 3",
			"poster":     "Poster 3",
			"is_watched": false,
			"status":     "downloading",
		},
		map[string]interface{}{
			"id":         "tt2",
			"rating":     9,
			"title":      "C'est le 2",
			"poster":     "Poster 2",
			"is_watched": true,
			"status":     "ready",
		},
		map[string]interface{}{
			"id":         "tt1",
			"rating":     7.5,
			"title":      "C'est le 1",
			"poster":     "Poster 1",
			"is_watched": false,
			"status":     "not downloaded",
		},
		map[string]interface{}{
			"id":         "tt9",
			"rating":     2.5,
			"title":      "This is 9",
			"poster":     "Poster 9",
			"is_watched": true,
			"status":     "ready",
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetCategoryOffsetNumberItems(t *testing.T) {
	tests.DbClean()
	category := "Testcategory"
	userID := "test_" + lib.GetRandomString(43)
	insertData(category, userID)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  userID,
	}
	r := tests.CreateRequest("GET", "/api/v1/movies/category/"+strings.ToLower(category)+"/1/2/it", []byte(`{}`), context)
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
			"id":         "tt2",
			"rating":     9,
			"title":      "This is 2",
			"poster":     "Poster 2",
			"is_watched": true,
			"status":     "ready",
		},
		map[string]interface{}{
			"id":         "tt1",
			"rating":     7.5,
			"title":      "è 1",
			"poster":     "Poster 1",
			"is_watched": false,
			"status":     "not downloaded",
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetCategoryOffsetNumberItemsNilUnknownLanguage(t *testing.T) {
	tests.DbClean()
	category := "Testcategory"
	userID := "test_" + lib.GetRandomString(43)
	insertData(category, userID)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  userID,
	}
	r := tests.CreateRequest("GET", "/api/v1/movies/category/"+strings.ToLower(category)+"/2/AZERT/pk", []byte(`{}`), context)
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
			"id":         "tt1",
			"rating":     7.5,
			"title":      "This is 1",
			"poster":     "Poster 1",
			"is_watched": false,
			"status":     "not downloaded",
		},
		map[string]interface{}{
			"id":         "tt9",
			"rating":     2.5,
			"title":      "This is 9",
			"poster":     "Poster 9",
			"is_watched": true,
			"status":     "ready",
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetCategoryUnknownCategory(t *testing.T) {
	tests.DbClean()
	category := "Testcategory"
	userID := "test_" + lib.GetRandomString(43)
	insertData(category, userID)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  userID,
	}
	r := tests.CreateRequest("GET", "/api/v1/movies/category/empty/0/20/en", []byte(`{}`), context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]string{
		"status": "No (more) data",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetCategoryPopular(t *testing.T) {
	tests.DbClean()
	category := "Testcategory"
	userID := "test_" + lib.GetRandomString(43)
	insertData(category, userID)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  userID,
	}
	r := tests.CreateRequest("GET", "/api/v1/movies/category/popular/0/20/fr", []byte(`{}`), context)
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
			"id":         "tt2",
			"rating":     9,
			"title":      "C'est le 2",
			"poster":     "Poster 2",
			"is_watched": true,
			"status":     "ready",
		},
		map[string]interface{}{
			"id":         "tt1",
			"rating":     7.5,
			"title":      "C'est le 1",
			"poster":     "Poster 1",
			"is_watched": false,
			"status":     "not downloaded",
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}
