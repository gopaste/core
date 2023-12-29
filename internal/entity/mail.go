package entity

import (
	"bytes"
	"html/template"
	"log"
	"os"

	"github.com/Caixetadev/snippet/internal/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

type MailData struct {
	Username string
	Code     string
}

func NewMailData(username string) MailData {
	return MailData{
		Username: username,
		Code:     utils.GenerateRandomString(8),
	}
}

const charset = "UTF-8"

func NewEmail(to string, htmlBody string, subject string, sourceEmail string) *ses.SendEmailInput {
	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(to)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sourceEmail),
	}
}

type EmailService interface {
	SendResetPasswordEmail(user *User) (string, error)
}

func GetHTMLTemplate(emailData MailData) string {
	var templateBuffer bytes.Buffer

	htmlData, err := os.ReadFile("./web/template/password_recovery.html")
	if err != nil {
		return ""
	}

	htmlTemplate := template.Must(template.New("password_recovery.html").Parse(string(htmlData)))

	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "password_recovery.html", emailData)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	return templateBuffer.String()
}
