package profile

import (
	"net/http/httptest"
	"testing"

	"../../../lib"
	coltype "../../../mongodb/collections"
	"../../../tests"
)

func TestGetUser(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(coltype.User{Username: "fake°1"}, tests.MongoDB)
	username := "test_" + lib.GetRandomString(43)
	_ = tests.InsertUser(coltype.User{
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
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("GET", "/v1/users/"+username, []byte(`{}`), context)
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
		"firstname":       "thisFirstname",
		"lastname":        "thisLastname",
		"profile_picture": "/storage/thisPicture",
		"username":        username,
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetUserInvalidUsername(t *testing.T) {
	tests.DbClean()
	username := "test_567898765(§è!çè§((§è!)))"
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("GET", "/v1/users/"+username, []byte(`{}`), context)
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
		"error": "Target username is invalid",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetUserNotFound(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(coltype.User{Username: "fake°1"}, tests.MongoDB)
	username := "test_" + lib.GetRandomString(43)
	_ = tests.InsertUser(coltype.User{Username: "fake°2"}, tests.MongoDB)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("GET", "/v1/users/"+username, []byte(`{}`), context)
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
}
