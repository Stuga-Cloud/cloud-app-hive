package services

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailService struct {
	host     string
	port     string
	username string
	password string
}

func NewEmailService() *EmailService {
	return &EmailService{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_PASSWORD"),
	}
}

func (s *EmailService) Send(to, subject, body string) error {
	from := s.username
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	fmt.Println("Sending email to", to, "with subject", subject, "and body", body, "from", from)

	err := smtp.SendMail(
		s.host+":"+s.port,
		smtp.PlainAuth(
			"",
			s.username,
			s.password,
			s.host,
		),
		from, []string{to}, []byte(msg),
	)

	if err != nil {
		return err
	}
	return nil
}
