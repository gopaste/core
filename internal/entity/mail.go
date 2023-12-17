package entity

import "github.com/Caixetadev/snippet/internal/utils"

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
