package auth

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"../../../lib"
	coltype "../../../mongodb/collections"
	"../../../tests"
)

func TestLoginNoBody(t *testing.T) {
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/login", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := lib.CaptureOutput(func() {
		Login(w, r)
	})
	// Check : Content stardard output
	mustContains := "Failed to decode body EOF"
	if !strings.Contains(output, mustContains) {
		t.Errorf("Must print an error containing '%s', not '%s'", mustContains, output)
	}
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Failed to decode body"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestLoginNoDatabase(t *testing.T) {
	r := httptest.NewRequest("POST", "/v1/account/register", nil)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Login(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 500
	expectedContent := "Problem with database connection"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestLoginEmptyFields(t *testing.T) {
	body := []byte(`{"username": "vomnes", "password": ""}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/login", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Login(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 403
	expectedContent := "Cannot have an empty field"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestLoginWrongUsername(t *testing.T) {
	body := []byte(`{"username": "vomnes", "password": "abcABC123"}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/login", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Login(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 403
	expectedContent := "Username or password incorrect"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(coltype.User{Username: "vomnes", Email: "valentin@g.com", LastName: "Omnes", FirstName: "Valentin", Password: "abcABC123", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	body := []byte(`{"username": "vomnes", "password": "abc"}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/login", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Login(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 403
	expectedContent := "Username or password incorrect"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

// abcABC123 -> $2a$10$pgek6WtdhtKmGXPWOOtEf.gsgtNXOkqr3pBjaCCa9il6XhRS7LAua
func TestLogin(t *testing.T) {
	tests.DbClean()
	u := tests.InsertUser(coltype.User{Username: "vomnes", Password: "$2a$10$pgek6WtdhtKmGXPWOOtEf.gsgtNXOkqr3pBjaCCa9il6XhRS7LAua", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	body := []byte(`{"uuid": "uuid-test", "username": "vomnes", "password": "abcABC123"}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/login", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Login(w, r)
	resp := w.Result()
	expectedCode := 200
	var response map[string]interface{}
	if err := tests.ChargeResponse(w, &response); err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != expectedCode {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m.\n%#v", expectedCode, resp.StatusCode, response)
	}
	settingsJWT := lib.DataJWT{
		Duration:       time.Hour * time.Duration(24*31),
		ISS:            "hypertube.com",
		Sub:            "",
		UserID:         u.ID,
		Username:       u.Username,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		ProfilePicture: u.ProfilePicture,
	}
	expectedJWT, err := lib.GenerateJWT(settingsJWT)
	if err != nil {
		t.Error("GenerateJWT - Fail to generate expected JWT - ", err)
	}
	if response["token"] != expectedJWT {
		t.Errorf("Response token is not correct - Header\nExpect: '%s'\nHas   : '%s'", expectedJWT, response["token"])
	}
}
