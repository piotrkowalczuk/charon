package service

import (
	"html/template"

	"github.com/go-soa/charon/lib/routing"
)

// TemplatesConfig ...
type TemplatesConfig struct {
	Path string `xml:"path"`
}

// WebTemplates ...
var WebTemplates *template.Template

// MailTemplates ...
var MailTemplates *template.Template

// InitTemplates ...
func InitTemplates(config TemplatesConfig, urlGenerator routing.URLGenerator) {
	getTemplatePath := func(path string) string {
		return config.Path + "/" + path
	}

	var err error
	// Website tempaltes
	webTemplates := template.New("")
	webTemplates.Funcs(template.FuncMap{
		"url":     routing.URLTemplateFunc(URLGenerator),
		"url_abs": routing.URLAbsTemplateFunc(URLGenerator),
	})
	webTemplates, err = webTemplates.ParseFiles(
		getTemplatePath("header.html"),
		getTemplatePath("footer.html"),
		getTemplatePath("registration/index.html"),
		getTemplatePath("registration/success.html"),
	)
	if err != nil {
		Logger.Fatal(err)
	}

	// Mail templates
	mailTemplates := template.New("")
	mailTemplates.Funcs(template.FuncMap{
		"url":     routing.URLTemplateFunc(URLGenerator),
		"url_abs": routing.URLAbsTemplateFunc(URLGenerator),
	})
	mailTemplates, err = mailTemplates.ParseGlob(getTemplatePath("mail/*"))
	if err != nil {
		Logger.Fatal(err)
	}

	WebTemplates = webTemplates
	MailTemplates = mailTemplates
}
