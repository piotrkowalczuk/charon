package service

import gomail "gopkg.in/gomail.v1"

// SMTP ...
type SMTP struct {
	m    *gomail.Mailer
	from string
}

// SendWelcomeMail ...
func (s *SMTP) SendWelcomeMail(to string, toUsername string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Welcome!")
	msg.SetBody("text/html", "Welcome <b>"+toUsername+"</b> in our site!")

	return s.m.Send(msg)
}

// NewSMTP ...
func NewSMTP(config MailerConfig) Mailer {
	return &SMTP{
		m:    gomail.NewMailer(config.Host, config.Username, config.Password, config.Port),
		from: config.From,
	}
}
