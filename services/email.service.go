package services

import (
	"fmt"
	"log"
	"os"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type EmailService struct {
	email     string
	username  string
	publicKey string
	secretKey string
}

func NewEmailService() *EmailService {
	return &EmailService{
		email:     os.Getenv("SMTP_EMAIL"),
		username:  os.Getenv("SMTP_USERNAME"),
		publicKey: os.Getenv("MAIL_JET_API_KEY"),
		secretKey: os.Getenv("MAIL_JET_SECRET_KEY"),
	}
}

func (s *EmailService) Send(to, subject, body string, htmlBody string) error {
	fromEmail := s.email
	fromUsername := s.username

	mailjetClient := mailjet.NewMailjetClient(s.publicKey, s.secretKey)
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: fromEmail,
				Name:  fromUsername,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: to,
				},
			},
			Subject:  subject,
			TextPart: body,
			HTMLPart: htmlBody,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("Sent email to %s with status %s (subject: %s)\n", res.ResultsV31[0].To[0].Email, res.ResultsV31[0].Status, subject)
	return nil
}
