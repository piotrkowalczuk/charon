package mail

import (
	"bytes"
	"html/template"
)

// Mail ...
type Mail struct {
	transport Transporter
	from      string
	templates *template.Template
}

// SendWelcomeMail ...
func (m *Mail) SendWelcomeMail(toMail string, toUsername string) error {
	mailHTMLBody := new(bytes.Buffer)
	mailTopic := new(bytes.Buffer)

	tplParams := &welcomeMailTplParams{
		Username: toUsername,
		Mail:     toMail,
	}

	err := m.templates.ExecuteTemplate(mailHTMLBody, tplNameForWelcome, tplParams)
	if err != nil {
		return err
	}

	err = m.templates.ExecuteTemplate(mailTopic, tplNameForWelcomeTopic, nil)
	if err != nil {
		return err
	}

	m.transport.Send(m.from, toMail, mailTopic.String(), "", mailHTMLBody.String())
	return nil
}

// NewMail ...
func NewMail(transport Transporter, from string, templates *template.Template) *Mail {
	return &Mail{
		transport: transport,
		from:      from,
		templates: templates,
	}
}
