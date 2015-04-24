package mail

import gomail "gopkg.in/gomail.v1"

// TransportSMTP ...
type TransportSMTP struct {
	m *gomail.Mailer
}

// Send ...
func (t *TransportSMTP) Send(from string, to string, topic string, textBody string, htmlBody string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", topic)

	if htmlBody != "" {
		msg.SetBody("text/html", htmlBody)
	}

	if textBody != "" {
		msg.SetBody("text/plain", textBody)
	}

	return t.m.Send(msg)
}

// NewTransportSMTP ...
func NewTransportSMTP(host string, username string, password string, port int) *TransportSMTP {
	return &TransportSMTP{
		m: gomail.NewMailer(host, username, password, port),
	}
}
