package mail

import (
	"bytes"
	"io"
)

const (
	confirmationMailerDirectory = "registration_confirmation"
)

// Sender ...
type Sender interface {
	Send(string, map[string]interface{}) error
}

type templateGetter interface {
	GetForMail(writer io.Writer, templateName string, params interface{}) error
}

// Mailer ...
type Mailer struct {
	from      string
	transport Transporter
	templates templateGetter
}

// NewMailer ...
func NewMailer(directory, from string, transport Transporter, templates templateGetter) (*Mailer, error) {
	return &Mailer{
		from:      from,
		templates: templates,
		transport: transport,
	}, nil
}

// NewConfirmationMailer ...
func NewConfirmationMailer(from string, transport Transporter, templates templateGetter) (*Mailer, error) {
	return NewMailer(confirmationMailerDirectory, from, transport, templates)
}

// Send ...
func (m Mailer) Send(to string, params map[string]interface{}) error {
	var err error
	var topic, html, plain bytes.Buffer

	err = m.templates.GetForMail(&topic, "registration_confirmation_topic", nil)
	if err != nil {
		return err
	}

	err = m.templates.GetForMail(&plain, "registration_confirmation_plain_body", params)
	if err != nil {
		return err
	}

	err = m.templates.GetForMail(&html, "registration_confirmation_html_body", params)
	if err != nil {
		return err
	}

	bodies := map[string]string{
		"text/html":  html.String(),
		"text/plain": plain.String(),
	}
	return m.transport.Send(m.from, to, topic.String(), bodies)
}
