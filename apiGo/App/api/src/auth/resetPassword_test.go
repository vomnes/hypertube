package auth

import (
	"net/http/httptest"
	"strings"
	"testing"

	"../../../lib"
	coltypes "../../../mongodb/collections"
	"../../../mongodb/query"
	"../../../tests"

	"gopkg.in/mgo.v2/bson"
)

func TestResetPasswordNoDatabase(t *testing.T) {
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", nil, tests.ContextData{})
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 500
	expectedContent := "Problem with database connection"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestResetPasswordNoBody(t *testing.T) {
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := lib.CaptureOutput(func() {
		ResetPassword(w, r)
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

func TestResetPasswordNoRandomToken(t *testing.T) {
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	body := []byte(`{"randomToken": "", "password": "abcABC123", "rePassword": "abcABC123"}`)
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "No field inside the body can be empty"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestResetPasswordNoPassword(t *testing.T) {
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	body := []byte(`{"randomToken": "myAwesomeToken", "password": "", "rePassword": "abcABC123"}`)
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "No field inside the body can be empty"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestResetPasswordNoRePassword(t *testing.T) {
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	body := []byte(`{"randomToken": "myAwesomeToken", "password": "abcABC123", "rePassword": ""}`)
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "No field inside the body can be empty"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestResetPasswordNoIdenticalPassword(t *testing.T) {
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	body := []byte(`{"randomToken": "myAwesomeToken", "password": "abcABC123", "rePassword": "abcABC"}`)
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Both password entered must be identical"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestResetPasswordNoAValidPassword(t *testing.T) {
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	body := []byte(`{"randomToken": "myAwesomeToken", "password": "abcABC", "rePassword": "abcABC"}`)
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid password"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestResetPasswordNoTokenInTheDB(t *testing.T) {
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	body := []byte(`{"randomToken": "myAwesomeToken", "password": "abcABC123", "rePassword": "abcABC123"}`)
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 400
	expectedContent := "Random token does not exists in the database"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestResetPassword(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(coltypes.User{Username: "vomnes", Password: "abc", RandomToken: "myAwesomeToken", AccountType: coltypes.AccountType{Level: 1}}, tests.MongoDB)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	newPassword := "NewABCabc123"
	body := []byte(`{"randomToken": "myAwesomeToken", "password": "` + newPassword + `", "rePassword": "` + newPassword + `"}`)
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 202
	expectedContent := ""
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
	var user coltypes.User
	user, err := query.FindUser(bson.M{"username": "vomnes"}, tests.MongoDB)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	if !strings.Contains(user.Password, "$2a$10$") {
		t.Error("Password in database has not been updated")
	}
	if user.RandomToken != "" {
		t.Error("RandomToken in database must be empty not equal to " + user.RandomToken)
	}
}
