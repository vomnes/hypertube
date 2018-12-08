package profile

import (
	"net/http/httptest"
	"testing"

	"../../../lib"
	coltype "../../../mongodb/collections"
	"../../../tests"
)

func TestGetProfile(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(coltype.User{Username: "fake°1"}, tests.MongoDB)
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(coltype.User{
		Username:       username,
		FirstName:      "thisFirstname",
		LastName:       "thisLastname",
		Email:          "test@test.co",
		ProfilePicture: "/storage/thisPicture",
		Locale:         "en",
		AccountType:    coltype.AccountType{Level: 1},
	}, tests.MongoDB)
	_ = tests.InsertUser(coltype.User{Username: "fake°2"}, tests.MongoDB)
	context := tests.ContextData{
		MongoDB:  tests.MongoDB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("GET", "/v1/profile/me", []byte(`{}`), context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetProfile(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{
		"available_locales": lib.LocaleAvailable,
		"email":             "test@test.co",
		"firstname":         "thisFirstname",
		"lastname":          "thisLastname",
		"locale":            "en",
		"profile_picture":   "/storage/thisPicture",
		"username":          username,
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetProfileOAuth(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(coltype.User{Username: "fake°1"}, tests.MongoDB)
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(coltype.User{
		Username:       username,
		FirstName:      "thisFirstname",
		LastName:       "thisLastname",
		Email:          "test@test.co",
		ProfilePicture: "/storage/thisPicture",
		Locale:         "en",
		AccountType: coltype.AccountType{
			Type:  "gplus",
			Level: 2,
		},
	}, tests.MongoDB)
	_ = tests.InsertUser(coltype.User{Username: "fake°2"}, tests.MongoDB)
	context := tests.ContextData{
		MongoDB:  tests.MongoDB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("GET", "/v1/profile/me", []byte(`{}`), context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetProfile(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{
		"available_locales": lib.LocaleAvailable,
		"firstname":         "thisFirstname",
		"lastname":          "thisLastname",
		"email":             "test@test.co",
		"locale":            "en",
		"profile_picture":   "/storage/thisPicture",
		"oauth":             "gplus",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetProfileNotFound(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(coltype.User{Username: "fake°1"}, tests.MongoDB)
	username := "test_" + lib.GetRandomString(43)
	_ = tests.InsertUser(coltype.User{Username: "fake°2"}, tests.MongoDB)
	context := tests.ContextData{
		MongoDB:  tests.MongoDB,
		Username: username,
		UserID:   "123",
	}
	r := tests.CreateRequest("GET", "/v1/profile/me", []byte(`{}`), context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetProfile(w, r)
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
}
