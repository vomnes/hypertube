package auth

import (
	"log"
	"net/http"

	"../../../lib"
	"../../../mongodb/query"

	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type dataBody struct {
	RandomToken   string `json:"randomToken"`
	NewPassword   string `json:"password"`
	NewRePassword string `json:"rePassword"`
}

func checkInputBody(inputData dataBody) (int, string) {
	if inputData.RandomToken == "" || inputData.NewPassword == "" ||
		inputData.NewRePassword == "" {
		return 406, "No field inside the body can be empty"
	}
	if inputData.NewPassword != inputData.NewRePassword {
		return 406, "Both password entered must be identical"
	}
	if !lib.IsValidPassword(inputData.NewPassword) {
		return 406, "Not a valid password"
	}
	return 0, ""
}

func getUserFromRandomToken(r *http.Request, db *mgo.Database, randomToken string) (int, string) {
	_, err := query.FindUser(bson.M{"random_token": randomToken, "account_type.level": 1}, db)
	if err != nil {
		if err.Error() == "Not Found" {
			return 400, "Random token does not exists in the database"
		}
		log.Println(lib.PrettyError(r.URL.String() + " [DB REQUEST - SELECT] " + err.Error()))
		return 500, "Failed to check if random token exists"
	}
	return 0, ""
}

func updateUserPassword(r *http.Request, db *mgo.Database, password, randomToken string) (int, string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(lib.PrettyError(r.URL.String() + " [PW - BCRYPT] " + err.Error()))
		return 500, "Password encryption failed"
	}
	err = query.UpdateUsers(bson.M{"random_token": randomToken, "account_type.level": 1}, bson.M{"password": string(hashedPassword), "random_token": ""}, db)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - UPDATE - RESETPASSWORD] Failed to update user data" + err.Error()))
		return 500, "UPDATE - Database request failed"
	}
	return 0, ""
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.MongoDB).(*mgo.Database)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	// Get body data
	var inputData dataBody
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
	errCode, errContent = getUserFromRandomToken(r, db, inputData.RandomToken)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = updateUserPassword(r, db, inputData.NewPassword, inputData.RandomToken)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondEmptyHTTP(w, http.StatusAccepted)
}
