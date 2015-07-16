package lib

import (
	"errors"

	gomail "gopkg.in/gomail.v1"
)

var (
	// ErrMissingMailBody ...
	ErrMissingMailBody = errors.New("charon: missing mail body")
)

const (
	// TransporterTypeSMTP ...
	TransporterTypeSMTP = "smtp"
)

// SMTPTransporter ...
type SMTPTransporter struct {
	m *gomail.Mailer
}

// Send ...
func (t *SMTPTransporter) Send(from, to, topic string, bodies map[string]string) error {
	msg := gomail.NewMessage()
	msg.SetHeaders(map[string][]string{
		"From":    {from},
		"To":      {to},
		"Subject": {topic},
	})

	if len(bodies) == 0 {
		return ErrMissingMailBody
	}

	for contentType, body := range bodies {
		msg.SetBody(contentType, body)
	}

	return t.m.Send(msg)
}

// NewSMTPTransporter ...
func NewSMTPTransporter(host string, username string, password string, port int) *SMTPTransporter {
	return &SMTPTransporter{
		m: gomail.NewMailer(host, username, password, port),
	}
}
