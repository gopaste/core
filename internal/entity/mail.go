package entity

import (
	"bytes"
	"html/template"
	"log"
	"os"

	"github.com/Caixetadev/snippet/internal/utils"
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
