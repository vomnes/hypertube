package mail

import (
	"log"
	"net/http"
	"os"

	"../../../lib"
	coltype "../../../mongodb/collections"
	query "../../../mongodb/query"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type userData struct {
	EmailAddress string `json:"email"`
	Test         bool   `json:"test"`
}

func checkEmailAddress(r *http.Request, db *mgo.Database, emailAddress string) (coltype.User, int, string) {
	var user coltype.User
	user, err := query.FindUser(bson.M{"email": emailAddress, "account_type.level": 1}, db)
	if err != nil {
		if err.Error() == "Not Found" {
			return coltype.User{}, 400, "Email address does not exists in the database"
		}
		log.Println(lib.PrettyError(r.URL.String() + " [DB REQUEST - SELECT] " + err.Error()))
		return coltype.User{}, 500, "Check if email address exists failed"
	}
	return user, 0, ""
}

func insertTokenDatabase(db *mgo.Database, emailAddress, randomToken string) (int, string) {
	err := query.UpdateUsers(bson.M{"email": emailAddress, "account_type.level": 1}, bson.M{"random_token": randomToken}, db)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - UPDATE - FORGOTPASSWORD] Failed to update user data with random_token" + err.Error()))
		return 500, "UPDATE - Database request failed"
	}
	return 0, ""
}

func sendMessage(w http.ResponseWriter, r *http.Request, isTest bool,
	user coltype.User, randomToken string, mailjetClient *mailjet.Client) {
	resetPasswordURL := os.Getenv("FRONT_DOMAIN_NAME") + "/recover/" + randomToken
	// Test response
	if isTest == true {
		lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"email":             user.Email,
			"fullname":          user.FirstName + " " + user.LastName,
			"forgotPasswordUrl": resetPasswordURL,
		})
		return
	}
	// Prod response
	mailVariables := map[string]interface{}{
		"firstname":         user.FirstName,
		"forgotPasswordUrl": resetPasswordURL,
	}
	subject := "Hypertube - Forgot password"
	if user.Locale == "fr" {
		subject = "Hypertube - Mot de passe oublié"
	} else if user.Locale == "it" {
		subject = "Hypertube - Password dimenticata"
	}
	templateID := lib.TemplateForgotPassword[user.Locale]
	if templateID == 0 {
		templateID = lib.TemplateForgotPassword["en"]
	}
	err := lib.SendMail(
		mailjetClient,
		user.Email,
		user.FirstName+" "+user.LastName,
		subject,
		templateID,
		mailVariables,
	)
	if err != nil {
		log.Println(lib.PrettyError(r.URL.String() + " [MAILJET - Forgot password] " + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Send forgot password email failed")
		return
	}
	lib.RespondEmptyHTTP(w, http.StatusAccepted)
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.MongoDB).(*mgo.Database)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	mailjetClient, ok := r.Context().Value(lib.MailJet).(*mailjet.Client)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with mailjet connection")
		return
	}
	// Get body data
	var inputData userData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	// Check email
	if inputData.EmailAddress == "" || !lib.IsValidEmailAddress(inputData.EmailAddress) {
		lib.RespondWithErrorHTTP(w, 400, "Email address is not valid")
		return
	}
	user, errCode, errContent := checkEmailAddress(r, db, inputData.EmailAddress)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	randomToken := lib.UniqueTimeToken(user.Username)
	// Insert random_token in database
	errCode, errContent = insertTokenDatabase(db, inputData.EmailAddress, randomToken)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	sendMessage(w, r, inputData.Test, user, randomToken, mailjetClient)
}
