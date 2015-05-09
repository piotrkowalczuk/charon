package mail

import (
	"bytes"
	"html/template"
)

const (
	confirmationMailerDirectory = "registration_confirmation"
)

// Sender ...
type Sender interface {
	Send(string, map[string]interface{}) error
}

// Mailer ...
type Mailer struct {
	from      string
	transport Transporter
	templates *template.Template
}

// NewMailer ...
func NewMailer(directory, from string, transport Transporter, templates *template.Template) (*Mailer, error) {
	templates, err := templates.ParseGlob(directory + "/*")

	if err != nil {
		return nil, err
	}

	return &Mailer{
		from:      from,
		templates: templates,
		transport: transport,
	}, nil
}

// NewConfirmationMailer ...
func NewConfirmationMailer(from string, transport Transporter, templates *template.Template) (*Mailer, error) {
	return NewMailer(confirmationMailerDirectory, from, transport, templates)
}

func (m Mailer) renderPlainBody(params map[string]interface{}) (string, error) {
	body := &bytes.Buffer{}
	err := m.templates.ExecuteTemplate(body, "/body.txt", params)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}

func (m Mailer) renderHTMLBody(params map[string]interface{}) (string, error) {
	body := &bytes.Buffer{}
	err := m.templates.ExecuteTemplate(body, "/body.html", params)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}

func (m Mailer) renderTitle() (string, error) {
	title := &bytes.Buffer{}
	err := m.templates.ExecuteTemplate(title, "/title.txt", nil)
	if err != nil {
		return "", err
	}

	return title.String(), nil
}

// Send ...
func (m Mailer) Send(to string, params map[string]interface{}) error {
	title, err := m.renderTitle()
	if err != nil {
		return err
	}
	plainBody, err := m.renderPlainBody(params)
	if err != nil {
		return err
	}
	htmlBody, err := m.renderHTMLBody(params)
	if err != nil {
		return err
	}

	bodies := map[string]string{
		"text/html":  htmlBody,
		"text/plain": plainBody,
	}
	return m.transport.Send(m.from, to, title, bodies)
}
