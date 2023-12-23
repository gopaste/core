package services

import (
	"fmt"

	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/Caixetadev/snippet/pkg/typesystem"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SimpleEmailService struct {
	sesClient *ses.SES
	Env       *config.Config
}

func NewSimpleEmailService(cfg *config.Config) (*SimpleEmailService, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AWSRegion),
		Credentials: credentials.NewStaticCredentials(cfg.AWSSecretKey, cfg.AWSAccessKey, ""),
	})

	if err != nil {
		return nil, err
	}

	return &SimpleEmailService{
		sesClient: ses.New(sess),
		Env:       cfg,
	}, nil
}

func (e *SimpleEmailService) SendResetPasswordEmail(user *entity.User) (string, error) {
	mailData := entity.NewMailData(user.Name)

	htmlBody := entity.GetHTMLTemplate(mailData)

	input := BuildEmail(user.Email, htmlBody, "Reset Password", e.Env.AWSSenderEmail)

	_, err := e.sesClient.SendEmail(input)
	if err != nil {
		fmt.Println(err)
		return "", typesystem.ServerError
	}

	return mailData.Code, nil
}

const charset = "UTF-8"

func BuildEmail(toEmail, htmlBody, subject, sourceEmail string) *ses.SendEmailInput {
	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(toEmail)},
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
