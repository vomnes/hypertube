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
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var invalidBodyFieldsPassword = []struct {
	body        map[string]interface{} // input
	errCode     int
	errStatus   string
	testContent string // test details
}{
	{
		map[string]interface{}{},
		406, "No field inside the body can be empty", "No field inside the body can be empty",
	},
	{
		map[string]interface{}{
			"password":       "xyz",
			"new_password":   "notNull",
			"new_rePassword": "notNull",
		}, 406, "Current password field is not a valid password", "Current password field is not a valid password",
	},
	{
		map[string]interface{}{
			"password":       "xyzXYZ123",
			"new_password":   "notEqul",
			"new_rePassword": "notEqual",
		}, 406, "Both password entered must be identical", "Both password entered must be identical",
	},
	{
		map[string]interface{}{
			"password":       "xyzXYZ123",
			"new_password":   "notValid",
			"new_rePassword": "notValid",
		}, 406, "New password field is not a valid password", "New password field is not a valid password",
	},
	{
		map[string]interface{}{
			"password":       "xyzXYZ123",
			"new_password":   "abcABC123",
			"new_rePassword": "abcABC123",
		}, 403, "Current password incorrect", "Current password incorrect",
	},
	{
		map[string]interface{}{
			"password":       "xyzXYZ123",
			"new_password":   "abcABC123",
			"new_rePassword": "abcABC123",
		}, 406, "User does not exists in the database", "User does not exists in the database",
	},
}

func TestEditPasswordCases(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	_ = tests.InsertUser(coltype.User{Username: "tester_1", Email: "random@gmail.com", Locale: "en", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	userData := tests.InsertUser(coltype.User{Username: username, Password: "fakePassword", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	for _, tt := range invalidBodyFieldsPassword {
		if tt.testContent == "User does not exists in the database" {
			username = "NotFoundUsername"
			userData.ID = "1234"
		}
		context := tests.ContextData{
			MongoDB:  tests.MongoDB,
			Username: username,
			UserID:   userData.ID,
		}
		body, err := lib.InterfaceToByte(tt.body)
		if err != nil {
			t.Error(err)
		}
		r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
		r.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		output := tests.CaptureOutput(func() {
			EditPassword(w, r)
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
		if tt.testContent != "User does not exists in the database" {
			user, err := query.FindUser(bson.M{"username": username, "account_type.level": 1}, tests.MongoDB)
			if err != nil {
				t.Errorf("Test type: %s\n%v", tt.testContent, err)
				return
			}
			expectedDatabase := coltype.User{
				ID:       userData.ID,
				Username: username,
				Password: "fakePassword",
				AccountType: coltype.AccountType{
					Level: 1,
				},
			}
			if compare := pretty.Compare(&expectedDatabase, user); compare != "" {
				t.Errorf("Test type: %s\n%s", tt.testContent, compare)
			}
		}
	}
}

func TestEditPasswordInvalidBody(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(coltype.User{Username: username}, tests.MongoDB)
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
		EditPassword(w, r)
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

// abcABC123 -> $2a$10$pgek6WtdhtKmGXPWOOtEf.gsgtNXOkqr3pBjaCCa9il6XhRS7LAua
func TestEditPassword(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(coltype.User{Username: username, Password: "$2a$10$pgek6WtdhtKmGXPWOOtEf.gsgtNXOkqr3pBjaCCa9il6XhRS7LAua", AccountType: coltype.AccountType{Level: 1}}, tests.MongoDB)
	context := tests.ContextData{
		MongoDB:  tests.MongoDB,
		Username: username,
		UserID:   userData.ID,
	}
	body, err := lib.InterfaceToByte(map[string]interface{}{
		"password":       "abcABC123",
		"new_password":   "abcABC1232019",
		"new_rePassword": "abcABC1232019",
	})
	if err != nil {
		t.Error(err)
	}
	r := tests.CreateRequest("POST", "/v1/profiles/edit/password", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		EditPassword(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	// Check : Updated data in database
	var user coltype.User
	user, err = query.FindUser(bson.M{"username": username, "account_type.level": 1}, tests.MongoDB)
	if err != nil {
		t.Error(err)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("abcABC1232019"))
	if err != nil {
		t.Error("New password not inserted in the database")
	}
}
