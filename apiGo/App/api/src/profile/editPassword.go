package profile

import (
	"log"
	"net/http"

	"../../../lib"
	coltype "../../../mongodb/collections"
	"../../../mongodb/query"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"golang.org/x/crypto/bcrypt"
)

type userPassword struct {
	CurrentPassword string `json:"password"`
	NewPassword     string `json:"new_password"`
	NewRePassword   string `json:"new_rePassword"`
}

func checkInputBody(inputData userPassword) (int, string) {
	if inputData.CurrentPassword == "" || inputData.NewPassword == "" ||
		inputData.NewRePassword == "" {
		return 406, "No field inside the body can be empty"
	}
	if !lib.IsValidPassword(inputData.CurrentPassword) {
		return 406, "Current password field is not a valid password"
	}
	if inputData.NewPassword != inputData.NewRePassword {
		return 406, "Both password entered must be identical"
	}
	if !lib.IsValidPassword(inputData.NewPassword) {
		return 406, "New password field is not a valid password"
	}
	return 0, ""
}

func checkCurrentUserPassword(r *http.Request, db *mgo.Database, password, userID, username string) (int, string) {
	var user coltype.User
	user, err := query.FindUser(bson.M{"_id": userID, "username": username, "account_type.level": 1}, db)
	if err != nil {
		if err.Error() == "Not Found" {
			return 406, "User does not exists in the database"
		}
		log.Println(lib.PrettyError(r.URL.String() + " [DB REQUEST - SELECT] " + err.Error()))
		return 500, "Failed to check if users exists in the database"
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 403, "Current password incorrect"
	}
	return 0, ""
}

func updateUserPassword(r *http.Request, db *mgo.Database, password, userID, username string) (int, string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(lib.PrettyError(r.URL.String() + " [PW - BCRYPT] " + err.Error()))
		return 500, "Password encryption failed"
	}
	err = query.UpdateUsers(
		bson.M{"_id": userID, "username": username},
		bson.M{"password": string(hashedPassword), "random_token": ""},
		db)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - UPDATE - EDITPASSWORD] Failed to update user data" + err.Error()))
		return 500, "UPDATE - Database request failed"
	}
	return 0, ""
}

// EditPassword is the route '/v1/profiles/edit/password' with the method POST.
// Return HTTP Code 200 Status OK
func EditPassword(w http.ResponseWriter, r *http.Request) {
	// Get bascis
	db, ok := r.Context().Value(lib.MongoDB).(*mgo.Database)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	username, _ := r.Context().Value(lib.Username).(string)
	userID, _ := r.Context().Value(lib.UserID).(string)
	var inputData userPassword
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = checkInputBody(inputData)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = checkCurrentUserPassword(r, db, inputData.CurrentPassword, userID, username)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = updateUserPassword(r, db, inputData.NewPassword, userID, username)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondEmptyHTTP(w, http.StatusOK)
}
