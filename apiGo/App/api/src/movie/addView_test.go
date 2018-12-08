package movie

import (
	"net/http/httptest"
	"testing"
	"time"

	"../../../lib"
	coltypes "../../../mongodb/collections"
	"../../../mongodb/query"
	"../../../tests"
	"github.com/kylelemons/godebug/pretty"
	"gopkg.in/mgo.v2/bson"
)

func TestAddView(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	arrayID := "test_" + lib.GetRandomString(43)
	user1 := tests.InsertUser(coltypes.User{FirstName: "firstname1", LastName: "1lastname", ProfilePicture: "picture1", Locale: "en"}, tests.MongoDB)
	_ = tests.InsertMovie(coltypes.Movie{
		IDimdb:        imdb,
		OriginalTitle: "This is 1",
		Poster:        "Poster 1",
		WatchedBy: []coltypes.Watched{
			coltypes.Watched{
				UserID:    arrayID,
				WatchedAt: time.Date(2018, 2, 2, 10, 1, 0, 0, time.UTC),
			},
		},
	}, tests.MongoDB)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  user1.ID,
	}
	tests.TimeTest = time.Now()
	r := tests.CreateRequest("POST", "/api/v1/movies/view/"+imdb, nil, context)
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
	dbState, err := query.FindMovies(bson.M{}, bson.M{"_id": 0}, tests.MongoDB)
	if err != nil {
		t.Errorf("%v", err)
	}
	expectedDatabaseState := []coltypes.Movie{
		coltypes.Movie{
			IDimdb:        "",
			OriginalTitle: "This is 1",
			Titles:        map[string]string{},
			Year:          0,
			Duration:      0,
			Genres:        []string{},
			Rating: coltypes.Rating{
				Average: 0,
				Number:  0,
			},
			Poster: "Poster 1",
			Video: coltypes.Video{
				Status:    0,
				Path:      "",
				Stream:    false,
				Subtitles: []coltypes.Subtitle{},
			},
			WatchedBy: []coltypes.Watched{
				coltypes.Watched{
					UserID:    arrayID,
					WatchedAt: time.Date(2018, 2, 2, 11, 1, 0, 0, time.Local),
				},
				coltypes.Watched{
					UserID:    user1.ID,
					WatchedAt: time.Now(),
				},
			},
		},
	}
	if compare := pretty.Compare(&expectedDatabaseState, dbState); compare != "" {
		t.Error(compare)
	}
}

func TestAddViewUnknowMovie(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	user1 := tests.InsertUser(coltypes.User{FirstName: "firstname1", LastName: "1lastname", ProfilePicture: "picture1", Locale: "en"}, tests.MongoDB)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  user1.ID,
	}
	tests.TimeTest = time.Now()
	r := tests.CreateRequest("POST", "/api/v1/movies/view/"+imdb, nil, context)
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

func TestAddViewNoUser(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	arrayID := "test_" + lib.GetRandomString(43)
	user1 := "user_" + lib.GetRandomString(43)
	_ = tests.InsertMovie(coltypes.Movie{
		IDimdb:        imdb,
		OriginalTitle: "This is 1",
		Poster:        "Poster 1",
		WatchedBy: []coltypes.Watched{
			coltypes.Watched{
				UserID:    arrayID,
				WatchedAt: time.Date(2018, 2, 2, 10, 1, 0, 0, time.UTC),
			},
		},
	}, tests.MongoDB)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
		UserID:  user1,
	}
	tests.TimeTest = time.Now()
	r := tests.CreateRequest("POST", "/api/v1/movies/view/"+imdb, nil, context)
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
		"error": "User does not exists in the database",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	dbState, err := query.FindMovies(bson.M{}, bson.M{"_id": 0}, tests.MongoDB)
	if err != nil {
		t.Errorf("%v", err)
	}
	expectedDatabaseState := []coltypes.Movie{
		coltypes.Movie{
			IDimdb:        "",
			OriginalTitle: "This is 1",
			Titles:        map[string]string{},
			Year:          0,
			Duration:      0,
			Genres:        []string{},
			Rating: coltypes.Rating{
				Average: 0,
				Number:  0,
			},
			Poster: "Poster 1",
			Video: coltypes.Video{
				Status:    0,
				Path:      "",
				Stream:    false,
				Subtitles: []coltypes.Subtitle{},
			},
			WatchedBy: []coltypes.Watched{
				coltypes.Watched{
					UserID:    arrayID,
					WatchedAt: time.Date(2018, 2, 2, 11, 1, 0, 0, time.Local),
				},
			},
		},
	}
	if compare := pretty.Compare(&expectedDatabaseState, dbState); compare != "" {
		t.Error(compare)
	}
}

func TestAddViewNoDB(t *testing.T) {
	tests.DbClean()
	imdb := "myFilmId"
	user1 := tests.InsertUser(coltypes.User{FirstName: "firstname1", LastName: "1lastname", ProfilePicture: "picture1", Locale: "en"}, tests.MongoDB)
	context := tests.ContextData{
		UserID: user1.ID,
	}
	tests.TimeTest = time.Now()
	r := tests.CreateRequest("POST", "/api/v1/movies/view/"+imdb, nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 500, map[string]interface{}{
		"error": "Problem with database connection",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}
