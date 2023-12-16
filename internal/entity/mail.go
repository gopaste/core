package entity

type MailData struct {
	Username string
	Code     string
}

type EmailService interface {
	SendResetPasswordEmail(to string, emailData MailData) error
}
