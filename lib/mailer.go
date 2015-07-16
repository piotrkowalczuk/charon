package lib

import (
	"bytes"
	"io"
)

// Sender ...
type Sender interface {
	Send(string, interface{}) error
}

type templateGetter interface {
	GetForMail(writer io.Writer, templateName string, params interface{}) error
}

// Mailer ...
type Mailer struct {
	directory string
	from      string
	transport Transporter
	templates templateGetter
}

// NewMailer ...
func NewMailer(directory, from string, transport Transporter, templates templateGetter) (*Mailer, error) {
	return &Mailer{
		directory: directory,
		from:      from,
		templates: templates,
		transport: transport,
	}, nil
}

// Send ...
func (m Mailer) Send(to string, params interface{}) error {
	var err error
	var topic, html, plain bytes.Buffer

	err = m.templates.GetForMail(&topic, m.directory+"_topic", nil)
	if err != nil {
		return err
	}

	err = m.templates.GetForMail(&plain, m.directory+"_plain_body", params)
	if err != nil {
		return err
	}

	err = m.templates.GetForMail(&html, m.directory+"_html_body", params)
	if err != nil {
		return err
	}

	bodies := map[string]string{
		"text/html":  html.String(),
		"text/plain": plain.String(),
	}
	return m.transport.Send(m.from, to, topic.String(), bodies)
}
