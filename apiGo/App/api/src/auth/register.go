package auth

import (
	"log"
	"net/http"
	"os"
	"strings"

	"../../../lib"
	coltypes "../../../mongodb/collections"
	query "../../../mongodb/query"
	"golang.org/x/crypto/bcrypt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type accountData struct {
	Username     string `json:"username"`
	EmailAddress string `json:"email"`
	Lastname     string `json:"lastname"`
	Firstname    string `json:"firstname"`
	Password     string `json:"password"`
	RePassword   string `json:"rePassword"`
	Base64       string `json:"picture_base64"`
}

func checkInput(d accountData) (int, string) {
	if d.Username == "" || d.EmailAddress == "" || d.Firstname == "" ||
		d.Lastname == "" || d.Password == "" || d.RePassword == "" {
		return 406, "At least one field of the body is empty"
	}
	right := lib.IsValidUsername(d.Username)
	if right == false {
		return 406, "Not a valid username"
	}
	right = lib.IsValidFirstLastName(d.Firstname)
	if right == false {
		return 406, "Not a valid firstname"
	}
	right = lib.IsValidFirstLastName(d.Lastname)
	if right == false {
		return 406, "Not a valid lastname"
	}
	right = lib.IsValidEmailAddress(d.EmailAddress)
	if right == false {
		return 406, "Not a valid email address"
	}
	if d.Password != d.RePassword {
		return 406, "Both password entered must be identical"
	}
	right = lib.IsValidPassword(d.Password)
	if right == false {
		return 406, "Not a valid password"
	}
	if d.Base64 == "" {
		return 406, "Base64 can't be empty"
	}
	return 0, ""
}

// availabilityInput check the validity in the database of the username and
// email address in order to avoid duplicates
func availabilityInput(d accountData, db *mgo.Database, r *http.Request) (int, string) {
	usernameInput := d.Username
	emailInput := d.EmailAddress
	var users []coltypes.User
	users, err := query.FindUsers(bson.M{"$or": []bson.M{bson.M{"username": usernameInput, "account_type.level": 1}, bson.M{"email": emailInput, "account_type.level": 1}}}, db)
	if err != nil {
		log.Println(lib.PrettyError(r.URL.String() + " [DB REQUEST - SELECT] " + err.Error()))
		return 406, "Check availability input failed"
	}
	usernameIsAvailable := true
	emailIsAvailable := true
	for _, user := range users {
		if user.Username == usernameInput {
			usernameIsAvailable = false
		}
		if user.Email == emailInput {
			emailIsAvailable = false
		}
		if !usernameIsAvailable && !emailIsAvailable {
			return 406, "Username and email address already used"
		}
	}
	if !usernameIsAvailable {
		return 406, "Username already used"
	} else if !emailIsAvailable {
		return 406, "Email address already used"
	}
	return 0, ""
}

func handlePicture(base64, username string) (string, int, string) {
	path, err := os.Getwd()
	if err != nil {
		return "", 500, "Failed to get the root path name - EncodeBase64"
	}
	path = strings.TrimSuffix(strings.TrimSuffix(path, "/src/auth"), "/api")
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

func createUser(d accountData, picturePath string, db *mgo.Database, r *http.Request) (int, string) {
	// Generate "hash" to store from user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(d.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(lib.PrettyError(r.URL.String() + " [PW - BCRYPT] " + err.Error()))
		return 500, "Password encryption failed"
	}
	// Insert data in database
	profile := coltypes.User{
		ID:             lib.GetRandomString(42),
		Username:       d.Username,
		FirstName:      d.Firstname,
		LastName:       d.Lastname,
		ProfilePicture: picturePath,
		Email:          d.EmailAddress,
		Password:       string(hashedPassword),
		AccountType: coltypes.AccountType{
			Level: 1,
			Type:  "password",
		},
		Locale: "en",
	}
	err = query.InsertUser(profile, db)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request insert user" + err.Error()))
		return 500, "Insert data failed"
	}
	return 0, ""
}

// Register is the route '/v1/account/register' with the method POST.
func Register(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.MongoDB).(*mgo.Database)
	if !ok {
		log.Println(lib.PrettyError("Register - Database Connection Failed"))
		lib.RespondWithErrorHTTP(w, http.StatusInternalServerError, "Problem with database connection")
		return
	}
	var inputData accountData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = checkInput(inputData)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = availabilityInput(inputData, db, r)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	picturePath, errCode, errContent := handlePicture(inputData.Base64, inputData.Username)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	picturePath = strings.TrimPrefix(picturePath, "/storage")
	errCode, errContent = createUser(inputData, picturePath, db, r)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondEmptyHTTP(w, http.StatusCreated)
}
