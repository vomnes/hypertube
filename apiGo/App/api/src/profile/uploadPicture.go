package profile

import (
	"log"
	"net/http"
	"os"
	"strings"

	"../../../lib"
	coltype "../../../mongodb/collections"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type pictureData struct {
	Base64 string `json:"picture_base64"`
}

func checkBase64(base64 string) (int, string) {
	if base64 == "" {
		return 406, "Base64 can't be empty"
	}
	if !lib.IsValidBase64Picture(base64) {
		return 406, "Base64 doesn't match with the pattern 'data:image/[...];base64,[...]'"
	}
	return 0, ""
}

func handlePicture(base64, path, username string) (string, int, string) {
	typeImage, imageBase64, err := lib.ExtractDataPictureBase64(base64)
	if err != nil {
		return "", 406, err.Error()
	}
	subPath := "/storage/pictures/profiles/"
	if username == "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7" {
		subPath = "/storage/tests/"
	}
	newPath, err := lib.Base64ToImageFile(imageBase64, typeImage, path, subPath, username)
	if err != nil {
		log.Println(lib.PrettyError("[PICTURE] " + err.Error()))
		return "", 406, lib.TrimStringFromString(err.Error(), " |")
	}
	return newPath, 0, ""
}

func updatePicturePath(r *http.Request, db *mgo.Database, picturePath, userID, username string) (string, int, string) {
	// Get oldPicturePath and replace path
	var user coltype.User
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"profile_picture": picturePath}},
		ReturnNew: false,
	}
	_, err := db.
		C("users").
		Find(bson.M{"_id": userID, "username": username}).
		Apply(change, &user)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - UPDATE - UPLOADPICTURE] Failed to update user data " + err.Error()))
		return "", 500, "UPDATE - Database request failed"
	}
	return user.ProfilePicture, 0, ""
}

func UploadPicture(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.MongoDB).(*mgo.Database)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	username, _ := r.Context().Value(lib.Username).(string)
	userID, _ := r.Context().Value(lib.UserID).(string)
	var inputData pictureData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = checkBase64(inputData.Base64)
	if errCode != 0 {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	path, err := os.Getwd()
	if err != nil {
		lib.RespondWithErrorHTTP(w, 500, "Failed to get the root path name")
		return
	}
	path = strings.TrimSuffix(strings.TrimSuffix(path, "/src/profile"), "/api")
	picturePath, errCode, errContent := handlePicture(inputData.Base64, path, username)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	picturePath = strings.TrimPrefix(picturePath, "/storage")
	oldPicturePath, errCode, errContent := updatePicturePath(r, db, picturePath, userID, username)
	if errCode != 0 {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	if oldPicturePath != "" && strings.Contains(oldPicturePath, "/storage/") {
		err = os.Remove(path + oldPicturePath)
		if err != nil {
			log.Println(lib.PrettyError("[OS] Failed to remove old picture - " + username + " - " + err.Error()))
		}
	}
	lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"picture_url": picturePath,
	})
}
