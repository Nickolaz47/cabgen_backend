package email

import (
	gomail "gopkg.in/mail.v2"
)

type EmailSender interface {
	Send(email *gomail.Message) error
}

type EmailConfig struct {
	Sender    string
	Recipient string
	Subject   string
	Body      string
	File      string
}

type SMTPEmailSender struct {
	Host     string
	Port     int
	Username string
	Password string
}

type EmailService struct {
	Config EmailConfig
	Sender EmailSender
}

func (s *SMTPEmailSender) Send(message *gomail.Message) error {
	dialer := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)
	return dialer.DialAndSend(message)
}

func setupEmailMessage(message *gomail.Message, emailConfig EmailConfig) {
	sender := emailConfig.Sender
	recipient := emailConfig.Recipient

	message.SetHeader("From", sender)
	message.SetHeader("To", recipient)
	message.SetHeader("Subject", emailConfig.Subject)
	message.SetBody("text/html", emailConfig.Body)

	if emailConfig.File != "" {
		message.Attach(emailConfig.File)
	}
}

func SendEmail(emailConfig EmailConfig, sender EmailSender) error {
	m := gomail.NewMessage()
	setupEmailMessage(m, emailConfig)

	return sender.Send(m)
}
