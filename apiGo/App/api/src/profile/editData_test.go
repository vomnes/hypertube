package profile

import (
	"net/http/httptest"
	"strings"
	"testing"

	"../../../lib"
	coltype "../../../mongodb/collections"
	"../../../mongodb/query"
	"../../../tests"
	"github.com/kylelemons/godebug/pretty"
	"gopkg.in/mgo.v2/bson"
)

func TestEditData(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(coltype.User{Username: username, Locale: "en", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	context := tests.ContextData{
		MongoDB:  tests.MongoDB,
		Username: username,
		UserID:   userData.ID,
	}
	body, err := lib.InterfaceToByte(map[string]interface{}{
		"firstname": "  valentin  ",
		"lastname":  "omnes  ",
		"email":     "valentin.omnes@gmail.com",
		"locale":    "fr",
	})
	if err != nil {
		t.Error(err)
	}
	r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		EditData(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	user, err := query.FindUser(bson.M{"username": username, "account_type.level": 1}, tests.MongoDB)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := coltype.User{
		ID:        userData.ID,
		Username:  username,
		Email:     "valentin.omnes@gmail.com",
		LastName:  "Omnes",
		FirstName: "Valentin",
		Locale:    "fr",
		AccountType: coltype.AccountType{
			Level: 1,
		},
	}
	if compare := pretty.Compare(&expectedDatabase, user); compare != "" {
		t.Error(compare)
	}
}

func TestEditDataEmailUsedByOAuth(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	_ = tests.InsertUser(coltype.User{Username: "something_random", Email: "valentin.omnes@gmail.com", AccountType: coltype.AccountType{Level: 2, Type: "gplus"}}, tests.MongoDB)
	userData := tests.InsertUser(coltype.User{Username: username, Email: "valentin.omnes@gmail.com", Locale: "en", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	context := tests.ContextData{
		MongoDB:  tests.MongoDB,
		Username: username,
		UserID:   userData.ID,
	}
	body, err := lib.InterfaceToByte(map[string]interface{}{
		"firstname": "  valentin  ",
		"lastname":  "omnes  ",
		"email":     "valentin.omnes@gmail.com",
		"locale":    "fr",
	})
	if err != nil {
		t.Error(err)
	}
	r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		EditData(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	user, err := query.FindUser(bson.M{"username": username, "account_type.level": 1}, tests.MongoDB)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := coltype.User{
		ID:        userData.ID,
		Username:  username,
		Email:     "valentin.omnes@gmail.com",
		LastName:  "Omnes",
		FirstName: "Valentin",
		Locale:    "fr",
		AccountType: coltype.AccountType{
			Level: 1,
		},
	}
	if compare := pretty.Compare(&expectedDatabase, user); compare != "" {
		t.Error(compare)
	}
}

func TestEditDataOAuth(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	_ = tests.InsertUser(coltype.User{Username: "tester_1", Email: "random@gmail.com", Locale: "en", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	userData := tests.InsertUser(coltype.User{Username: username, Email: "v@g.co", Locale: "en", LastName: "TesterY", FirstName: "TesterX", AccountType: coltype.AccountType{Level: 2, Type: "gplus"}}, tests.MongoDB)
	context := tests.ContextData{
		MongoDB:  tests.MongoDB,
		Username: username,
		UserID:   userData.ID,
	}
	body, err := lib.InterfaceToByte(map[string]interface{}{
		"email":    "fake@test.com",
		"lastname": "tester",
		"locale":   "ru",
	})
	if err != nil {
		t.Error(err)
	}
	r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		EditData(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	user, err := query.FindUser(bson.M{"username": username, "account_type.level": 2}, tests.MongoDB)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := coltype.User{
		ID:        userData.ID,
		Username:  username,
		Email:     "v@g.co",
		FirstName: "TesterX",
		LastName:  "Tester",
		Locale:    "en",
		AccountType: coltype.AccountType{
			Level: 2,
			Type:  "gplus",
		},
	}
	if compare := pretty.Compare(&expectedDatabase, user); compare != "" {
		t.Error(compare)
	}
}

var invalidFields = []struct {
	body        map[string]interface{} // input
	errCode     int
	errStatus   string
	testContent string // test details
}{
	{map[string]interface{}{}, 400, "Nothing to update", "Nothing to update"},
	{map[string]interface{}{"firstname": "@@@@@rdftagzhdjkazd"}, 406, "Not a valid firstname", "Not a valid firstname"},
	{map[string]interface{}{"lastname": "@"}, 406, "Not a valid lastname", "Not a valid lastname"},
	{map[string]interface{}{"email": "12345678%1"}, 406, "Not a valid email address", "Not a valid email address"},
	{map[string]interface{}{"email": "random@gmail.com"}, 406, "Email address already used by an other user", "Email address already used by an other user"},
}

func TestEditDataCases(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	_ = tests.InsertUser(coltype.User{Username: "tester_1", Email: "random@gmail.com", Locale: "en", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	userData := tests.InsertUser(coltype.User{Username: username, Email: "v@test.co", Locale: "en", LastName: "TesterY", FirstName: "TesterX", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	context := tests.ContextData{
		MongoDB:  tests.MongoDB,
		Username: username,
		UserID:   userData.ID,
	}
	for _, tt := range invalidFields {
		body, err := lib.InterfaceToByte(tt.body)
		if err != nil {
			t.Error(err)
		}
		r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
		r.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		output := tests.CaptureOutput(func() {
			EditData(w, r)
		})
		// Check : Content stardard output
		if output != "" {
			t.Errorf("Test type: %s\n%s", tt.testContent, output)
		}
		strError := tests.CompareResponseJSONCode(w, tt.errCode, map[string]interface{}{
			"error": tt.errStatus,
		})
		if strError != nil {
			t.Errorf("Test type: %s\n%v", tt.testContent, strError)
		}
		user, err := query.FindUser(bson.M{"username": username, "account_type.level": 1}, tests.MongoDB)
		if err != nil {
			t.Errorf("Test type: %s\n%v", tt.testContent, err)
			return
		}
		expectedDatabase := coltype.User{
			ID:        userData.ID,
			Username:  username,
			Email:     "v@test.co",
			LastName:  "TesterY",
			FirstName: "TesterX",
			Locale:    "en",
			AccountType: coltype.AccountType{
				Level: 1,
			},
		}
		if compare := pretty.Compare(&expectedDatabase, user); compare != "" {
			t.Errorf("Test type: %s\n%s", tt.testContent, compare)
		}
	}
}

func TestEditDataInvalidBody(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	_ = tests.InsertUser(coltype.User{Username: "tester_1", Email: "random@gmail.com", Locale: "en", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	userData := tests.InsertUser(coltype.User{Username: username, Email: "v@g.co", Locale: "en", LastName: "TesterY", FirstName: "TesterX", AccountType: coltype.AccountType{Level: 2, Type: "gplus"}}, tests.MongoDB)
	context := tests.ContextData{
		MongoDB:  tests.MongoDB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		EditData(w, r)
	})
	// Check : Content stardard output
	mustContains := "Failed to decode body unexpected EOF"
	if !strings.Contains(output, mustContains) {
		t.Errorf("Must print an error containing '%s', not '%s'", mustContains, output)
	}
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Failed to decode body",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestEditDataUnknownUser(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		MongoDB:  tests.MongoDB,
		Username: "username",
		UserID:   "userData.ID",
	}
	body, err := lib.InterfaceToByte(map[string]interface{}{
		"firstname": "  valentin  ",
		"lastname":  "omnes  ",
		"email":     "valentin.omnes@gmail.com",
		"locale":    "fr",
	})
	if err != nil {
		t.Error(err)
	}
	r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		EditData(w, r)
	})
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 500, map[string]interface{}{
		"error": "Collect user data failed",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}
