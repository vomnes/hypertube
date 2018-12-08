package auth

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"../../../lib"
	coltype "../../../mongodb/collections"
	"../../../mongodb/query"
	"../../../tests"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func TestRegisterNoBody(t *testing.T) {
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := lib.CaptureOutput(func() {
		Register(w, r)
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

func TestRegisterFieldEmptyBody(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "rePassword": "abcABC123"}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "At least one field of the body is empty"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldInvalidUsername(t *testing.T) {
	body := []byte(`{"username": "vomnes&&", "email": "vomnes@student.42.fr", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "rePassword": "abcABC123"}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid username"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldInvalidFirstname(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "vomnes@student.42.fr", "firstname": "Valentin..", "lastname": "Omnes", "password": "abcABC123", "rePassword": "abcABC123"}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid firstname"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldInvalidLastname(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "vomnes@student.42.fr", "firstname": "Valentin", "lastname": "Omnes**", "password": "abcABC123", "rePassword": "abcABC123"}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid lastname"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldInvalidEmailAddress(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "vomnes", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "rePassword": "abcABC123"}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid email address"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldWrongIndenticalPassword(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "vomnes@s.co", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "rePassword": "abcABC123Wrong"}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Both password entered must be identical"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldInvalidPassword(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "vomnes@s.co", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC", "rePassword": "abcABC"}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid password"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterNotAvailableUsernameEmail(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(coltype.User{Username: "vomn", Email: "valentin@gmail.com", LastName: "Omnes", FirstName: "Valentin", Password: "abc", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	_ = tests.InsertUser(coltype.User{Username: "vomnes", Email: "valentin.omnes@gmail.com", LastName: "Omnes", FirstName: "Valentin", Password: "abc", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	body := []byte(`{"username": "vomnes", "email": "valentin.omnes@gmail.com", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "rePassword": "abcABC123", "picture_base64": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQIAIwAjAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAALCAAQABABAREA/8QAFwAAAwEAAAAAAAAAAAAAAAAAAQMECP/EAB4QAAEFAQEBAQEAAAAAAAAAAAQDBQYHCAIJARQS/9oACAEBAAA/AN17njdi6FtO8qT3puJvJpfL26WBopqu+cL1XZkPn5r55wzbXDKRaMeOloAcgb4lX7jeMdZhHxBxAPmtfQ6Qp/mfyRfgIwnHbEz/AGtR9J4L3E3i0np7ccnbLorzrC9V1nEIAfHfOeCaweFqtjoMuPDYG2YQQSjY69isaLaA3zKwJlIlPhMgFJ4Pv9PeacC3lqas7e1HnPKhku5xtqCCO2lJ8vXjFYQRuR/Q7Fs9HhhiLO8qvLxGjZDXar6lwMkgEzFfelSv19AAmt8ukqbL3fmOtKh1HnLVZkRG2Xp+du2a58vYbFXwY+VPOfGMDFmZizOzKszzJTI3YqzEj2Mqgayif2iV+vhwBC//2Q=="}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Username and email address already used"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterNotAvailableUsername(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(coltype.User{Username: "vomnvv", Email: "valentin@g.com", LastName: "Omnes", FirstName: "Valentin", Password: "abc", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	_ = tests.InsertUser(coltype.User{Username: "vomnes", Email: "valentin@gmail.com", LastName: "Omnes", FirstName: "Valentin", Password: "abc", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	body := []byte(`{"username": "vomnes", "email": "valentin.omnes@gmail.com", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "rePassword": "abcABC123", "picture_base64": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQIAIwAjAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAALCAAQABABAREA/8QAFwAAAwEAAAAAAAAAAAAAAAAAAQMECP/EAB4QAAEFAQEBAQEAAAAAAAAAAAQDBQYHCAIJARQS/9oACAEBAAA/AN17njdi6FtO8qT3puJvJpfL26WBopqu+cL1XZkPn5r55wzbXDKRaMeOloAcgb4lX7jeMdZhHxBxAPmtfQ6Qp/mfyRfgIwnHbEz/AGtR9J4L3E3i0np7ccnbLorzrC9V1nEIAfHfOeCaweFqtjoMuPDYG2YQQSjY69isaLaA3zKwJlIlPhMgFJ4Pv9PeacC3lqas7e1HnPKhku5xtqCCO2lJ8vXjFYQRuR/Q7Fs9HhhiLO8qvLxGjZDXar6lwMkgEzFfelSv19AAmt8ukqbL3fmOtKh1HnLVZkRG2Xp+du2a58vYbFXwY+VPOfGMDFmZizOzKszzJTI3YqzEj2Mqgayif2iV+vhwBC//2Q=="}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Username already used"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterNotAvailableEmailAddress(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(coltype.User{Username: "vomnvv", Email: "valentin@g.com", LastName: "Omnes", FirstName: "Valentin", Password: "abc", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	_ = tests.InsertUser(coltype.User{Username: "vomnvv", Email: "valentin.omnes@gmail.com", LastName: "Omnes", FirstName: "Valentin", Password: "abc", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	body := []byte(`{"username": "vomnes", "email": "valentin.omnes@gmail.com", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "rePassword": "abcABC123", "picture_base64": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQIAIwAjAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAALCAAQABABAREA/8QAFwAAAwEAAAAAAAAAAAAAAAAAAQMECP/EAB4QAAEFAQEBAQEAAAAAAAAAAAQDBQYHCAIJARQS/9oACAEBAAA/AN17njdi6FtO8qT3puJvJpfL26WBopqu+cL1XZkPn5r55wzbXDKRaMeOloAcgb4lX7jeMdZhHxBxAPmtfQ6Qp/mfyRfgIwnHbEz/AGtR9J4L3E3i0np7ccnbLorzrC9V1nEIAfHfOeCaweFqtjoMuPDYG2YQQSjY69isaLaA3zKwJlIlPhMgFJ4Pv9PeacC3lqas7e1HnPKhku5xtqCCO2lJ8vXjFYQRuR/Q7Fs9HhhiLO8qvLxGjZDXar6lwMkgEzFfelSv19AAmt8ukqbL3fmOtKh1HnLVZkRG2Xp+du2a58vYbFXwY+VPOfGMDFmZizOzKszzJTI3YqzEj2Mqgayif2iV+vhwBC//2Q=="}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Email address already used"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
	var u coltype.User
	u, err := query.FindUser(bson.M{"username": "vomnes", "email": "valentin.omnes@gmail.com"}, tests.MongoDB)
	if err != nil && err.Error() != "Not Found" {
		t.Error("Error with mongoDB in FindUser")
		return
	}
	if err != nil && (u.Username != "" || u.Email != "") {
		t.Error("User must hasn't been inserted", err)
	}
}

func TestRegisterEmptyBase64(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(coltype.User{Username: "vomnvv", Email: "valentin@g.com", LastName: "Omnes", FirstName: "Valentin", Password: "abc", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	_ = tests.InsertUser(coltype.User{Username: "vomnvv", Email: "valentin.omnes@gmail.com", LastName: "Omnes", FirstName: "Valentin", Password: "abc", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	body := []byte(`{"username": "vomnes", "email": "valentin.omnes@gmail.com", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "rePassword": "abcABC123", "picture_base64": ""}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Base64 can't be empty"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
	var u coltype.User
	u, err := query.FindUser(bson.M{"username": "vomnes", "email": "valentin.omnes@gmail.com"}, tests.MongoDB)
	if err != nil && err.Error() != "Not Found" {
		t.Error("Error with mongoDB in FindUser")
		return
	}
	if err != nil && (u.Username != "" || u.Email != "") {
		t.Error("User must not has been inserted", err)
	}
}

func TestRegisterInvalidBase64(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(coltype.User{Username: "vomnvv", Email: "valentin@g.com", LastName: "Omnes", FirstName: "Valentin", Password: "abc", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	_ = tests.InsertUser(coltype.User{Username: "vomnvv", Email: "valentin.omnes@gmail.com", LastName: "Omnes", FirstName: "Valentin", Password: "abc", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	body := []byte(`{"username": "vomnes", "email": "valentin.omnes@g.com", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "rePassword": "abcABC123", "picture_base64": "data:image/jpeg;,azdazd"}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Base64 doesn't match with the pattern 'data:image/[...];base64,[...]'"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
	var u coltype.User
	u, err := query.FindUser(bson.M{"username": "vomnes", "email": "valentin.omnes@gmail.com"}, tests.MongoDB)
	if err != nil && err.Error() != "Not Found" {
		t.Error("Error with mongoDB in FindUser")
		return
	}
	if err != nil && (u.Username != "" || u.Email != "") {
		t.Error("User must not has been inserted", err)
	}
}

func TestRegisterCorruptedPicture(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(coltype.User{Username: "vomnvv", Email: "valentin@g.com", LastName: "Omnes", FirstName: "Valentin", Password: "abc", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	_ = tests.InsertUser(coltype.User{Username: "vomnvv", Email: "valentin.omnes@gmail.com", LastName: "Omnes", FirstName: "Valentin", Password: "abc", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	body := []byte(`{"username": "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7", "email": "valentin.omnes@g.com", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "rePassword": "abcABC123", "picture_base64": "data:image/jpeg;base64,"}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := lib.CaptureOutput(func() {
		Register(w, r)
	})
	// Check : Content stardard output
	mustContains := "[PICTURE] Corrupted file [image/jpeg] | Error:"
	if !strings.Contains(output, mustContains) {
		t.Errorf("Must print an error containing '%s', not '%s'", mustContains, output)
	}
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Corrupted file [image/jpeg]"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
	var u coltype.User
	u, err := query.FindUser(bson.M{"username": "vomnes", "email": "valentin.omnes@gmail.com"}, tests.MongoDB)
	if err != nil && err.Error() != "Not Found" {
		t.Error("Error with mongoDB in FindUser")
		return
	}
	if err != nil && (u.Username != "" || u.Email != "") {
		t.Error("User must not has been inserted", err)
	}
}

func TestRegisterNoDatabase(t *testing.T) {
	r := httptest.NewRequest("POST", "/v1/account/register", nil)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := lib.CaptureOutput(func() {
		Register(w, r)
	})
	// Check : Content stardard output
	mustContains := "Register - Database Connection Failed"
	if !strings.Contains(output, mustContains) {
		t.Errorf("Must print an error containing '%s', not '%s'", mustContains, output)
	}
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 500
	expectedContent := "Problem with database connection"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegister(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(coltype.User{Username: "valentin", Email: "valentin@g.com", LastName: "Omnes", FirstName: "Valentin", Password: "abc", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	_ = tests.InsertUser(coltype.User{Username: "vomnes", Email: "valentin.omnes@gmail.com", LastName: "Omnes", FirstName: "Valentin", Password: "abc", AccountType: coltype.AccountType{Level: 2}}, tests.MongoDB)
	body := []byte(`{"username": "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7", "email": "valentin.omnes@gmail.com", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "rePassword": "abcABC123", "picture_base64": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQIAIwAjAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAALCAAQABABAREA/8QAFwAAAwEAAAAAAAAAAAAAAAAAAQMECP/EAB4QAAEFAQEBAQEAAAAAAAAAAAQDBQYHCAIJARQS/9oACAEBAAA/AN17njdi6FtO8qT3puJvJpfL26WBopqu+cL1XZkPn5r55wzbXDKRaMeOloAcgb4lX7jeMdZhHxBxAPmtfQ6Qp/mfyRfgIwnHbEz/AGtR9J4L3E3i0np7ccnbLorzrC9V1nEIAfHfOeCaweFqtjoMuPDYG2YQQSjY69isaLaA3zKwJlIlPhMgFJ4Pv9PeacC3lqas7e1HnPKhku5xtqCCO2lJ8vXjFYQRuR/Q7Fs9HhhiLO8qvLxGjZDXar6lwMkgEzFfelSv19AAmt8ukqbL3fmOtKh1HnLVZkRG2Xp+du2a58vYbFXwY+VPOfGMDFmZizOzKszzJTI3YqzEj2Mqgayif2iV+vhwBC//2Q=="}`)
	context := tests.ContextData{
		MongoDB: tests.MongoDB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	expectedCode := 201
	if resp.StatusCode != expectedCode {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m.", expectedCode, resp.StatusCode)
	}
	var u coltype.User
	u, err := query.FindUser(bson.M{"username": "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7", "email": "valentin.omnes@gmail.com", "account_type.level": 1}, tests.MongoDB)
	if err != nil && err.Error() != "Not Found" {
		t.Error("Error with mongoDB in FindUser")
	}
	if u.Username != "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7" {
		t.Errorf("Username stored in the database is not correct, expect \x1b[1;32m%s\033[0m has \x1b[1;31m%s\033[0m.", "vomnes", u.Username)
	}
	if u.Email != "valentin.omnes@gmail.com" {
		t.Errorf("Email address stored in the database is not correct, expect \x1b[1;32m%s\033[0m has \x1b[1;31m%s\033[0m.", "valentin.omnes@gmail.com", u.Email)
	}
	if u.FirstName != "Valentin" {
		t.Errorf("Firstname stored in the database is not correct, expect \x1b[1;32m%s\033[0m has \x1b[1;31m%s\033[0m.", "Omnes", u.FirstName)
	}
	if u.LastName != "Omnes" {
		t.Errorf("Lastname stored in the database is not correct, expect \x1b[1;32m%s\033[0m has \x1b[1;31m%s\033[0m.", "Valentin", u.LastName)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("abcABC123")); err != nil {
		t.Error("Password stored in the database is not correct")
	}
	if !strings.Contains(u.ProfilePicture, "/tests/") {
		t.Error("Picture path stored in the database is not correct contains '/tests/' but", u.ProfilePicture)
	}
	if u.Locale != "en" {
		t.Errorf("Locale stored in the database must be set to 'en' not to '%s'", u.Locale)
	}
	// Check : File created
	path, err := os.Getwd()
	if err != nil {
		t.Error("Failed to get the root path name - EncodeBase64")
	}
	path = strings.TrimSuffix(strings.TrimSuffix(path, "/src/auth"), "/api")
	empty, err := lib.FileExists(path + "/storage" + u.ProfilePicture)
	if err != nil {
		t.Error(err)
	}
	if !empty {
		t.Errorf("The file picture doesn't exists, it hasn't been created")
	}
	err = os.RemoveAll(path + "/storage/tests/")
	if err != nil {
		t.Error(err)
	}
}
