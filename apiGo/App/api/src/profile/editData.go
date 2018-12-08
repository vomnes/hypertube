package profile

import (
	"html"
	"log"
	"net/http"
	"strings"

	"../../../lib"
	"../../../mongodb/query"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type userData struct {
	EmailAddress string `json:"email"`
	LastName     string `json:"lastname"`
	FirstName    string `json:"firstname"`
	Locale       string `json:"locale"`
}

func checkDataInput(d *userData) (int, string) {
	if d.FirstName == "" && d.LastName == "" && d.EmailAddress == "" &&
		d.Locale == "" {
		return 400, "Nothing to update"
	}
	if d.FirstName != "" {
		d.FirstName = strings.Trim(d.FirstName, " ")
		d.FirstName = html.EscapeString(d.FirstName)
		if !lib.IsValidFirstLastName(d.FirstName) {
			return 406, "Not a valid firstname"
		}
		d.FirstName = strings.Title(d.FirstName)
	}
	if d.LastName != "" {
		d.LastName = strings.Trim(d.LastName, " ")
		d.LastName = html.EscapeString(d.LastName)
		if !lib.IsValidFirstLastName(d.LastName) {
			return 406, "Not a valid lastname"
		}
		d.LastName = strings.Title(d.LastName)
	}
	if d.EmailAddress != "" {
		d.EmailAddress = strings.Trim(d.EmailAddress, " ")
		d.EmailAddress = html.EscapeString(d.EmailAddress)
		if !lib.IsValidEmailAddress(d.EmailAddress) {
			return 406, "Not a valid email address"
		}
		d.EmailAddress = strings.ToLower(d.EmailAddress)
	}
	if d.Locale != "" {
		if !lib.StringInArray(d.Locale, lib.LocaleAvailable) {
			d.Locale = "en"
		}
	}
	return 0, ""
}

func checkEmailAddressAvailability(db *mgo.Database, username string, d *userData) (int, string) {
	_, err := query.FindUser(bson.M{"email": (*d).EmailAddress, "account_type.level": 1, "username": bson.M{"$ne": username}}, db)
	if err != nil {
		if err.Error() == "Not Found" {
			return 0, ""
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Check email address availability in the database " + err.Error()))
		return 500, "Failed to check if users exists in the database"
	}
	return 406, "Email address already used by an other user"
}

func updateDataInDB(db *mgo.Database, data userData, userID, username string) (int, string, error) {
	change := make(map[string]interface{})
	if data.EmailAddress != "" {
		change["email"] = data.EmailAddress
	}
	if data.FirstName != "" {
		change["firstname"] = data.FirstName
	}
	if data.LastName != "" {
		change["lastname"] = data.LastName
	}
	if data.Locale != "" {
		change["locale"] = data.Locale
	}
	err := query.UpdateUsers(bson.M{"_id": userID, "username": username}, change, db)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - Update] Failed to update User[" + userID + "] Profile Data " + err.Error()))
		return 500, "Failed to update data in database", err
	}
	return 0, "", nil
}

func EditData(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.MongoDB).(*mgo.Database)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	username, _ := r.Context().Value(lib.Username).(string)
	userID, _ := r.Context().Value(lib.UserID).(string)
	var inputData userData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	user, err := query.FindUser(bson.M{"_id": userID}, db)
	if err != nil {
		lib.RespondWithErrorHTTP(w, 500, "Collect user data failed")
		return
	}
	// Clear necessary fields for OAuth profiles
	if user.AccountType.Level == 2 {
		inputData.EmailAddress = ""
	}
	errCode, errContent = checkDataInput(&inputData)
	if errCode != 0 && errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	if inputData.EmailAddress != "" {
		errCode, errContent = checkEmailAddressAvailability(db, username, &inputData)
		if errCode != 0 && errContent != "" {
			lib.RespondWithErrorHTTP(w, errCode, errContent)
			return
		}
	}
	errCode, errContent, err = updateDataInDB(db, inputData, userID, username)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondEmptyHTTP(w, http.StatusOK)
}
