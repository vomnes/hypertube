package movie

import (
	"encoding/base64"
	"net/http/httptest"
	"testing"
	"time"

	"../../../lib"
	coltypes "../../../mongodb/collections"
	"../../../tests"
)

func insertMovies(category, userID string) {
	_ = tests.InsertMovie(coltypes.Movie{
		IDimdb:        "tt1",
		OriginalTitle: "This is 1",
		Titles: map[string]string{
			"it": "è 1",
			"fr": "C'est le 1",
		},
		Year:     2015,
		Duration: 150,
		Genres:   []string{category, "a", "b", "c"},
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
				WatchedAt: time.Date(2017, 2, 3, 10, 1, 0, 0, time.UTC),
			},
		},
	}, tests.MongoDB)
	_ = tests.InsertMovie(coltypes.Movie{
		IDimdb:        "tt9",
		OriginalTitle: "This is 9",
		Year:          2015,
		Duration:      60,
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
		OriginalTitle: "this is 2",
		Titles: map[string]string{
			"fr": "IS True",
		},
		Year:     2015,
		Duration: 75,
		Genres:   []string{"x", "y", category, "c"},
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
		Duration:      95,
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
		Duration:      45,
		Genres:        []string{"a"},
		Rating: coltypes.Rating{
			Average: 1,
		},
		Poster: "Poster 4",
		Video: coltypes.Video{
			Status: 0,
		},
	}, tests.MongoDB)
	_ = tests.InsertMovie(coltypes.Movie{
		IDimdb:        "tt5",
		OriginalTitle: "Emp.y",
		Year:          2014,
		Duration:      35,
		Genres:        []string{"a"},
		Rating: coltypes.Rating{
			Average: 8,
		},
		Poster: "Poster 5",
		Video: coltypes.Video{
			Status: 2,
		},
	}, tests.MongoDB)
	_ = tests.InsertMovie(coltypes.Movie{
		IDimdb:        "tt7",
		OriginalTitle: "Title x1",
		Year:          1915,
		Duration:      100,
		Genres:        []string{"a"},
		Rating: coltypes.Rating{
			Average: 8,
		},
		Poster: "Poster x1",
		Video: coltypes.Video{
			Status: 0,
		},
	}, tests.MongoDB)
	_ = tests.InsertMovie(coltypes.Movie{
		IDimdb:        "tt8",
		OriginalTitle: "Title x2",
		Year:          1925,
		Duration:      125,
		Genres:        []string{"a", "b", "c"},
		Rating: coltypes.Rating{
			Average: 7.5,
		},
		Poster: "Poster x2",
		Video: coltypes.Video{
			Status: 1,
		},
	}, tests.MongoDB)
	_ = tests.InsertMovie(coltypes.Movie{
		IDimdb:        "ttx9",
		OriginalTitle: "Title x3",
		Year:          1889,
		Duration:      49,
		Genres:        []string{"a"},
		Rating: coltypes.Rating{
			Average: 5,
		},
		Poster: "Poster x3",
		Video: coltypes.Video{
			Status: 2,
		},
	}, tests.MongoDB)
	_ = tests.InsertMovie(coltypes.Movie{
		IDimdb:        "tt10",
		OriginalTitle: "Title x4",
		Year:          1850,
		Duration:      39,
		Genres:        []string{"a", "y", "z"},
		Rating: coltypes.Rating{
			Average: 6.5,
		},
		Poster: "Poster x4",
		Video: coltypes.Video{
			Status: 2,
		},
	}, tests.MongoDB)
	_ = tests.InsertMovie(coltypes.Movie{
		IDimdb:        "tt11",
		OriginalTitle: "Title x5",
		Year:          1945,
		Duration:      0,
		Genres:        []string{"a"},
		Rating: coltypes.Rating{
			Average: 4.25,
		},
		Poster: "Poster x5",
		Video: coltypes.Video{
			Status: -1,
		},
	}, tests.MongoDB)
	_ = tests.InsertMovie(coltypes.Movie{
		IDimdb:        "tt12",
		OriginalTitle: "Title x6",
		Year:          1900,
		Duration:      150,
		Genres:        []string{"a", "qwd"},
		Rating: coltypes.Rating{
			Average: 4.99,
		},
		Poster: "Poster x6",
		Video: coltypes.Video{
			Status: -1,
		},
	}, tests.MongoDB)
}

func TestSearchMovies(t *testing.T) {
	tests.DbClean()
	category := "Testcategory"
	userID := "test_" + lib.GetRandomString(43)
	insertMovies(category, userID)
	searchParameters, err := lib.InterfaceToByte(map[string]interface{}{
		"search": "t",
		"year": map[string]float32{
			"max": 2016,
			"min": 2000,
		},
		"rating": map[string]float32{
			"max": 9,
			"min": 0,
		},
	})
	if err != nil {
		t.Error(err)
	}
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  userID,
	}
	r := tests.CreateRequest("GET", "/api/v1/movies/search/0/20/en", nil, context)
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
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
			"title":      "this is 2",
			"poster":     "Poster 2",
			"is_watched": true,
			"status":     "ready",
			"year":       2015,
			"genres":     []string{"x", "y", "Testcategory", "c"},
			"duration":   75,
		},
		map[string]interface{}{
			"id":         "tt1",
			"rating":     7.5,
			"title":      "This is 1",
			"poster":     "Poster 1",
			"is_watched": false,
			"status":     "not downloaded",
			"year":       2015,
			"genres":     []string{"Testcategory", "a", "b", "c"},
			"duration":   150,
		},
		map[string]interface{}{
			"id":         "tt9",
			"rating":     2.5,
			"title":      "This is 9",
			"poster":     "Poster 9",
			"is_watched": true,
			"status":     "ready",
			"year":       2015,
			"genres":     []string{"x", "Testcategory", "y"},
			"duration":   60,
		},
		map[string]interface{}{
			"id":         "tt3",
			"rating":     1,
			"title":      "This is 3",
			"poster":     "Poster 3",
			"is_watched": false,
			"status":     "downloading",
			"year":       2016,
			"genres":     []string{"Testcategory"},
			"duration":   95,
		},
		map[string]interface{}{
			"id":         "tt4",
			"rating":     1,
			"title":      "This is 4",
			"poster":     "Poster 4",
			"is_watched": false,
			"status":     "waiting",
			"year":       2015,
			"genres":     []string{"a"},
			"duration":   45,
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestSearchMoviesLanguage(t *testing.T) {
	tests.DbClean()
	category := "Testcategory"
	userID := "test_" + lib.GetRandomString(43)
	insertMovies(category, userID)
	searchParameters, err := lib.InterfaceToByte(map[string]interface{}{
		"search": "is",
		"year": map[string]float32{
			"max": 2018,
			"min": 2000,
		},
		"rating": map[string]float32{
			"max": 9,
			"min": 5.5,
		},
	})
	if err != nil {
		t.Error(err)
	}
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  userID,
	}
	r := tests.CreateRequest("GET", "/api/v1/movies/search/0/abc/fr", nil, context)
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
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
			"title":      "IS True",
			"poster":     "Poster 2",
			"is_watched": true,
			"status":     "ready",
			"year":       2015,
			"genres":     []string{"x", "y", "Testcategory", "c"},
			"duration":   75,
		},
		map[string]interface{}{
			"id":         "tt1",
			"rating":     7.5,
			"title":      "C'est le 1",
			"poster":     "Poster 1",
			"is_watched": false,
			"status":     "not downloaded",
			"year":       2015,
			"genres":     []string{"Testcategory", "a", "b", "c"},
			"duration":   150,
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestSearchMoviesEmptyResponse(t *testing.T) {
	tests.DbClean()
	category := "Testcategory"
	userID := "test_" + lib.GetRandomString(43)
	insertMovies(category, userID)
	searchParameters, err := lib.InterfaceToByte(map[string]interface{}{
		"search": "is",
		"year": map[string]float32{
			"max": 2018,
			"min": 2000,
		},
		"rating": map[string]float32{
			"max": 9,
			"min": 5.5,
		},
	})
	if err != nil {
		t.Error(err)
	}
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  userID,
	}
	r := tests.CreateRequest("GET", "/api/v1/movies/search/20/20/fr", nil, context)
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
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
		"status": "No (more) data",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestSearchMoviesNegativeBody(t *testing.T) {
	tests.DbClean()
	category := "Testcategory"
	userID := "test_" + lib.GetRandomString(43)
	insertMovies(category, userID)
	searchParameters, err := lib.InterfaceToByte(map[string]interface{}{
		"search": "is",
		"year": map[string]float32{
			"max": -5,
			"min": 2000,
		},
		"rating": map[string]float32{
			"max": 9,
			"min": -2,
		},
	})
	if err != nil {
		t.Error(err)
	}
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  userID,
	}
	r := tests.CreateRequest("GET", "/api/v1/movies/search/20/20/fr", nil, context)
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
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
		"status": "No (more) data",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestSearchMoviesInvalidbase64(t *testing.T) {
	tests.DbClean()
	category := "Testcategory"
	userID := "test_" + lib.GetRandomString(43)
	insertMovies(category, userID)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  userID,
	}
	r := tests.CreateRequest("GET", "/api/v1/movies/search/20/20/fr", nil, context)
	r.Header.Add("Search-Parameters", "Hypertube")
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
		"error": "Failed to extract base64 search parameters in header",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestSearchMoviesDownloading(t *testing.T) {
	tests.DbClean()
	category := "Testcategory"
	userID := "test_" + lib.GetRandomString(43)
	insertMovies(category, userID)
	searchParameters, err := lib.InterfaceToByte(map[string]interface{}{
		"search": "x",
		"year": map[string]float32{
			"max": 1950,
			"min": 0,
		},
		"rating": map[string]float32{
			"max": 9,
			"min": 0,
		},
		"status": "downloading",
	})
	if err != nil {
		t.Error(err)
	}
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  userID,
	}
	r := tests.CreateRequest("GET", "/api/v1/movies/search/0/20/fr", nil, context)
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, []interface{}{
		map[string]interface{}{
			"id":         "tt7",
			"title":      "Title x1",
			"poster":     "Poster x1",
			"is_watched": false,
			"status":     "waiting",
			"rating":     8,
			"year":       1915,
			"genres":     []string{"a"},
			"duration":   100,
		},
		map[string]interface{}{
			"id":         "tt8",
			"title":      "Title x2",
			"poster":     "Poster x2",
			"is_watched": false,
			"status":     "downloading",
			"rating":     7.5,
			"year":       1925,
			"genres":     []string{"a", "b", "c"},
			"duration":   125,
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestSearchMoviesDownloaded(t *testing.T) {
	tests.DbClean()
	category := "Testcategory"
	userID := "test_" + lib.GetRandomString(43)
	insertMovies(category, userID)
	searchParameters, err := lib.InterfaceToByte(map[string]interface{}{
		"search": "x",
		"year": map[string]float32{
			"max": 1950,
			"min": 0,
		},
		"rating": map[string]float32{
			"max": 9,
			"min": 0,
		},
		"status": "ready",
	})
	if err != nil {
		t.Error(err)
	}
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  userID,
	}
	r := tests.CreateRequest("GET", "/api/v1/movies/search/0/20/en", nil, context)
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, []interface{}{
		map[string]interface{}{
			"id":         "tt10",
			"title":      "Title x4",
			"poster":     "Poster x4",
			"is_watched": false,
			"status":     "ready",
			"rating":     6.5,
			"year":       1850,
			"genres":     []string{"a", "y", "z"},
			"duration":   39,
		},
		map[string]interface{}{
			"id":         "ttx9",
			"title":      "Title x3",
			"poster":     "Poster x3",
			"is_watched": false,
			"status":     "ready",
			"rating":     5,
			"year":       1889,
			"genres":     []string{"a"},
			"duration":   49,
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestSearchMoviesNotDownloaded(t *testing.T) {
	tests.DbClean()
	category := "Testcategory"
	userID := "test_" + lib.GetRandomString(43)
	insertMovies(category, userID)
	searchParameters, err := lib.InterfaceToByte(map[string]interface{}{
		"search": "x",
		"year": map[string]float32{
			"max": 1950,
			"min": 0,
		},
		"rating": map[string]float32{
			"max": 9,
			"min": 0,
		},
		"status": "not downloaded",
	})
	if err != nil {
		t.Error(err)
	}
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  userID,
	}
	r := tests.CreateRequest("GET", "/api/v1/movies/search/0/20/fr", nil, context)
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, []interface{}{
		map[string]interface{}{
			"id":         "tt12",
			"title":      "Title x6",
			"poster":     "Poster x6",
			"is_watched": false,
			"status":     "not downloaded",
			"rating":     4.99,
			"year":       1900,
			"genres":     []string{"a", "qwd"},
			"duration":   150,
		},
		map[string]interface{}{
			"id":         "tt11",
			"title":      "Title x5",
			"poster":     "Poster x5",
			"is_watched": false,
			"status":     "not downloaded",
			"rating":     4.25,
			"year":       1945,
			"genres":     []string{"a"},
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestSearchMoviesAll(t *testing.T) {
	tests.DbClean()
	category := "Testcategory"
	userID := "test_" + lib.GetRandomString(43)
	insertMovies(category, userID)
	searchParameters, err := lib.InterfaceToByte(map[string]interface{}{
		"search": "x",
		"year": map[string]float32{
			"max": 1950,
			"min": 0,
		},
		"rating": map[string]float32{
			"max": 9,
			"min": 0,
		},
	})
	if err != nil {
		t.Error(err)
	}
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  userID,
	}
	r := tests.CreateRequest("GET", "/api/v1/movies/search/0/20/fr", nil, context)
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, []interface{}{
		map[string]interface{}{
			"id":         "tt7",
			"title":      "Title x1",
			"poster":     "Poster x1",
			"is_watched": false,
			"status":     "waiting",
			"rating":     8,
			"year":       1915,
			"genres":     []string{"a"},
			"duration":   100,
		},
		map[string]interface{}{
			"id":         "tt8",
			"title":      "Title x2",
			"poster":     "Poster x2",
			"is_watched": false,
			"status":     "downloading",
			"rating":     7.5,
			"year":       1925,
			"genres":     []string{"a", "b", "c"},
			"duration":   125,
		},
		map[string]interface{}{
			"id":         "tt10",
			"title":      "Title x4",
			"poster":     "Poster x4",
			"is_watched": false,
			"status":     "ready",
			"rating":     6.5,
			"year":       1850,
			"genres":     []string{"a", "y", "z"},
			"duration":   39,
		},
		map[string]interface{}{
			"id":         "ttx9",
			"title":      "Title x3",
			"poster":     "Poster x3",
			"is_watched": false,
			"status":     "ready",
			"rating":     5,
			"year":       1889,
			"genres":     []string{"a"},
			"duration":   49,
		},
		map[string]interface{}{
			"id":         "tt12",
			"title":      "Title x6",
			"poster":     "Poster x6",
			"is_watched": false,
			"status":     "not downloaded",
			"rating":     4.99,
			"year":       1900,
			"genres":     []string{"a", "qwd"},
			"duration":   150,
		},
		map[string]interface{}{
			"id":         "tt11",
			"title":      "Title x5",
			"poster":     "Poster x5",
			"is_watched": false,
			"status":     "not downloaded",
			"rating":     4.25,
			"year":       1945,
			"genres":     []string{"a"},
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestSearchMoviesEmptyParameters(t *testing.T) {
	tests.DbClean()
	category := "Testcategory"
	userID := "test_" + lib.GetRandomString(43)
	insertMovies(category, userID)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  userID,
	}
	r := tests.CreateRequest("GET", "/api/v1/movies/search/0/20/fr", nil, context)
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
		"status": "No (more) data",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}
