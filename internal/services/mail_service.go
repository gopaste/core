package services

import (
	"bytes"
	"html/template"
	"log"
	"os"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SimpleEmailService struct {
	sesClient *ses.SES
}

func NewSimpleEmailService() (*SimpleEmailService, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("AKIA2KTHNMQ7JQVJIQWA", "6uQhJjqYL0OcrBQ5Sr22ddVBhSC1+4QAIEXcNUKL", ""),
	})

	if err != nil {
		return nil, err
	}

	return &SimpleEmailService{
		sesClient: ses.New(sess),
	}, nil
}

func getHTMLTemplate(emailData entity.MailData) string {
	var templateBuffer bytes.Buffer

	htmlData, err := os.ReadFile("index.html")

	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))

	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", emailData)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	return templateBuffer.String()
}

func (e *SimpleEmailService) SendResetPasswordEmail(to string, emailData entity.MailData) error {
	html := getHTMLTemplate(emailData)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(to)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(html),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("Reset Password"),
			},
		},
		Source: aws.String("caixetadev@gmail.com"),
	}

	_, err := e.sesClient.SendEmail(input)

	return err
}
