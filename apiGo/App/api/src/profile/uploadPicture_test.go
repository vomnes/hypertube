package profile

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"../../../lib"
	coltype "../../../mongodb/collections"
	"../../../mongodb/query"
	"../../../tests"
	"github.com/kylelemons/godebug/pretty"
	"gopkg.in/mgo.v2/bson"
)

func getPathNameTest(t *testing.T) string {
	path, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	return strings.TrimSuffix(path, "/api/src/profile")
}

func TestUploadPicture(t *testing.T) {
	tests.DbClean()
	username := "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7"
	path := getPathNameTest(t)
	oldPicturePath := "/storage/tests/" + "thisIsTheUrl_TestPictureDelete"
	userData := tests.InsertUser(coltype.User{Username: username, ProfilePicture: oldPicturePath}, tests.MongoDB)
	os.MkdirAll(path+"/storage/tests/", os.ModePerm)
	f, err := os.Create(path + oldPicturePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	body, err := lib.InterfaceToByte(map[string]string{
		"picture_base64": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQIAIwAjAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAALCAAQABABAREA/8QAFwAAAwEAAAAAAAAAAAAAAAAAAQMECP/EAB4QAAEFAQEBAQEAAAAAAAAAAAQDBQYHCAIJARQS/9oACAEBAAA/AN17njdi6FtO8qT3puJvJpfL26WBopqu+cL1XZkPn5r55wzbXDKRaMeOloAcgb4lX7jeMdZhHxBxAPmtfQ6Qp/mfyRfgIwnHbEz/AGtR9J4L3E3i0np7ccnbLorzrC9V1nEIAfHfOeCaweFqtjoMuPDYG2YQQSjY69isaLaA3zKwJlIlPhMgFJ4Pv9PeacC3lqas7e1HnPKhku5xtqCCO2lJ8vXjFYQRuR/Q7Fs9HhhiLO8qvLxGjZDXar6lwMkgEzFfelSv19AAmt8ukqbL3fmOtKh1HnLVZkRG2Xp+du2a58vYbFXwY+VPOfGMDFmZizOzKszzJTI3YqzEj2Mqgayif2iV+vhwBC//2Q==",
	})
	if err != nil {
		t.Error(err)
	}
	context := tests.ContextData{
		MongoDB:  tests.MongoDB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("POST", "/v1/profiles/edit/picture", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	UploadPicture(w, r)
	user, err := query.FindUser(bson.M{"_id": userData.ID}, tests.MongoDB)
	if err != nil {
		t.Error(err)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{
		"picture_url": user.ProfilePicture,
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	// Check : File created
	empty, err := lib.FileExists(path + "/storage" + user.ProfilePicture)
	if err != nil {
		t.Error(err)
		return
	}
	if !empty {
		t.Errorf("The file picture doesn't exists, it hasn't been created")
	}
	// Check : Old file deleted
	empty, err = lib.FileExists(path + oldPicturePath)
	if err != nil {
		t.Error(err)
		return
	}
	if empty {
		t.Errorf("The old file picture exists, it hasn't been deleted")
	}
}

var invalidUploadPicture = []struct {
	body        map[string]interface{} // input
	errCode     int
	errStatus   string
	testContent string // test details
}{
	{
		map[string]interface{}{
			"picture_base64": "",
		}, 406, "Base64 can't be empty", "Base64 can't be empty",
	},
	{
		map[string]interface{}{
			"picture_base64": "data:image/jpeg;,/9j/",
		}, 406, "Base64 doesn't match with the pattern 'data:image/[...];base64,[...]'", "Base64 doesn't match with the pattern 'data:image/[...];base64,[...]'",
	},
	{
		map[string]interface{}{
			"picture_base64": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQIAIwAjAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEB",
		}, 406, "Corrupted file [image/jpeg]", "Corrupted file",
	},
}

func TestUploadPictureCases(t *testing.T) {
	tests.DbClean()
	username := "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7"
	userData := tests.InsertUser(coltype.User{Username: username, ProfilePicture: "empty"}, tests.MongoDB)
	context := tests.ContextData{
		MongoDB:  tests.MongoDB,
		Username: username,
		UserID:   userData.ID,
	}
	for _, tt := range invalidUploadPicture {
		body, err := lib.InterfaceToByte(tt.body)
		if err != nil {
			t.Error(err)
		}
		r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
		r.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		output := tests.CaptureOutput(func() {
			UploadPicture(w, r)
		})
		// Check : Content stardard output
		if tt.testContent != "Corrupted file" {
			if output != "" {
				t.Errorf("Test type: %s\n%s", tt.testContent, output)
			}
		} else {
			if !strings.Contains(output, "[PICTURE] Corrupted file [image/jpeg] | Error: unexpected EOF") {
				t.Errorf("Test type: %s\n%s", tt.testContent, output)
			}
		}
		strError := tests.CompareResponseJSONCode(w, tt.errCode, map[string]interface{}{
			"error": tt.errStatus,
		})
		if strError != nil {
			t.Errorf("Test type: %s\n%v", tt.testContent, strError)
		}
		user, err := query.FindUser(bson.M{"username": username}, tests.MongoDB)
		if err != nil {
			t.Errorf("Test type: %s\n%v", tt.testContent, err)
			return
		}
		expectedDatabase := coltype.User{
			ID:             userData.ID,
			Username:       username,
			ProfilePicture: "empty",
		}
		if compare := pretty.Compare(&expectedDatabase, user); compare != "" {
			t.Errorf("Test type: %s\n%s", tt.testContent, compare)
		}
	}
}

func TestUploadPictureInvalidBody(t *testing.T) {
	tests.DbClean()
	username := "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7"
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
		UploadPicture(w, r)
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

func TestUploadPictureNoOldPicture(t *testing.T) {
	tests.DbClean()
	username := "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7"
	path := getPathNameTest(t)
	userData := tests.InsertUser(coltype.User{Username: username, ProfilePicture: "", Email: "v@v.co", FirstName: "myFirstname", LastName: "myLastname"}, tests.MongoDB)
	body, err := lib.InterfaceToByte(map[string]string{
		"picture_base64": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQIAIwAjAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAALCAAQABABAREA/8QAFwAAAwEAAAAAAAAAAAAAAAAAAQMECP/EAB4QAAEFAQEBAQEAAAAAAAAAAAQDBQYHCAIJARQS/9oACAEBAAA/AN17njdi6FtO8qT3puJvJpfL26WBopqu+cL1XZkPn5r55wzbXDKRaMeOloAcgb4lX7jeMdZhHxBxAPmtfQ6Qp/mfyRfgIwnHbEz/AGtR9J4L3E3i0np7ccnbLorzrC9V1nEIAfHfOeCaweFqtjoMuPDYG2YQQSjY69isaLaA3zKwJlIlPhMgFJ4Pv9PeacC3lqas7e1HnPKhku5xtqCCO2lJ8vXjFYQRuR/Q7Fs9HhhiLO8qvLxGjZDXar6lwMkgEzFfelSv19AAmt8ukqbL3fmOtKh1HnLVZkRG2Xp+du2a58vYbFXwY+VPOfGMDFmZizOzKszzJTI3YqzEj2Mqgayif2iV+vhwBC//2Q==",
	})
	if err != nil {
		t.Error(err)
	}
	context := tests.ContextData{
		MongoDB:  tests.MongoDB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("POST", "/v1/profiles/edit/picture", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	UploadPicture(w, r)
	user, err := query.FindUser(bson.M{"_id": userData.ID}, tests.MongoDB)
	if err != nil {
		t.Error(err)
	}
	stateDB := coltype.User{
		ID:             user.ID,
		Username:       username,
		ProfilePicture: user.ProfilePicture,
		Email:          "v@v.co",
		FirstName:      "myFirstname",
		LastName:       "myLastname",
	}
	if compare := pretty.Compare(&stateDB, user); compare != "" {
		t.Error(compare)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{
		"picture_url": user.ProfilePicture,
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	// Check : File created
	empty, err := lib.FileExists(path + "/storage" + user.ProfilePicture)
	if err != nil {
		t.Error(err)
		return
	}
	if !empty {
		t.Errorf("The file picture doesn't exists, it hasn't been created")
	}
	err = os.RemoveAll(path + "/storage/tests/")
	if err != nil {
		t.Error(err)
	}
}
