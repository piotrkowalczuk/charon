package mail

import gomail "gopkg.in/gomail.v1"

// SMTPTransporter ...
type SMTPTransporter struct {
	m *gomail.Mailer
}

// Send ...
func (t *SMTPTransporter) Send(from, to, topic, contentType, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeaders(map[string][]string{
		"From":    {from},
		"To":      {to},
		"Subject": {topic},
	})
	msg.SetBody(contentType, body)

	return t.m.Send(msg)
}

// NewSMTPTransporter ...
func NewSMTPTransporter(host string, username string, password string, port int) *SMTPTransporter {
	return &SMTPTransporter{
		m: gomail.NewMailer(host, username, password, port),
	}
}
