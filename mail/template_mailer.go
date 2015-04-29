package mail

import (
	"bytes"
	"html/template"
	"log"
)

const (
	confirmationMailerTemplateName = "registration_confirmation"
)

// TemplateMailer ...
type TemplateMailer struct {
	mailer       Mailer
	templateName string

	templates *template.Template
}

// NewTemplateMailer ...
func NewTemplateMailer(mailer Mailer, templateName string, templates *template.Template) TemplateMailer {
	return TemplateMailer{
		mailer:       mailer,
		templateName: templateName,
		templates:    templates,
	}
}

// NewConfirmationMailer ...
func NewConfirmationMailer(mailer Mailer, templates *template.Template) TemplateMailer {
	return NewTemplateMailer(mailer, confirmationMailerTemplateName, templates)
}

// BodyFileName ...
func (tm TemplateMailer) BodyTemplateName() string {
	return tm.templateName + "_body"
}

// TitleFileName ...
func (tm TemplateMailer) TitleTemplateName() string {
	return tm.templateName + "_title"
}

// Send ...
func (tm TemplateMailer) Send(to string, params map[string]interface{}) error {
	mailHTMLBody := &bytes.Buffer{}
	mailTopic := &bytes.Buffer{}

	err := tm.templates.ExecuteTemplate(mailHTMLBody, tm.BodyTemplateName(), params)
	if err != nil {
		return err
	}

	err = tm.templates.ExecuteTemplate(mailTopic, tm.TitleTemplateName(), nil)
	if err != nil {
		return err
	}

	log.Println(mailTopic.String(), "text/html", mailHTMLBody.String(), tm.BodyTemplateName())

	return tm.mailer.transport.Send(tm.mailer.from, to, mailTopic.String(), "text/html", mailHTMLBody.String())
}
