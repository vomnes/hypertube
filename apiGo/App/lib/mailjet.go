package lib

import (
	"os"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

var (
	// TemplateForgotPassword corresponds to the id of the ForgotPassword template used in mailjet
	TemplateForgotPassword = map[string]int{
		"en": 582012,
		"fr": 606948,
		"it": 606949,
	}
)

// MailJetConn allows to create a connection with our MailJet client
func MailJetConn() *mailjet.Client {
	return mailjet.NewMailjetClient(os.Getenv("MJ_APIKEY_PUBLIC"), os.Getenv("MJ_APIKEY_PRIVATE"))
}

// SendMail allows to send email through mailjet
// Take as parameter mailjetClient, email and destinatary name, sibject
// mailjet templateId and variables containing the variables used in
// the template
func SendMail(mailjetClient *mailjet.Client, toEmail, toName, subject string,
	templateID int, variables map[string]interface{}) error {
	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: "valentin.omnes@gmail.com",
				Name:  "Mail Delivery - Hypertube",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: toEmail,
					Name:  toName,
				},
			},
			TemplateID:       templateID,
			TemplateLanguage: true,
			Subject:          subject,
			Variables:        variables,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		return err
	}
	return nil
}
