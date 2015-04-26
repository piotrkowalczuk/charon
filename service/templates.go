package service

import "html/template"

// TemplatesConfig ...
type TemplatesConfig struct {
	Path string `xml:"path"`
}

// Templates ...
var Templates *template.Template

// InitTemplates ...
func InitTemplates(config TemplatesConfig) {
	getTemplatePath := func(path string) string {
		return config.Path + "/" + path
	}

	templates, err := template.ParseFiles(
		getTemplatePath("header.html"),
		getTemplatePath("footer.html"),
		getTemplatePath("registration/index.html"),
		getTemplatePath("registration/success.html"),
	)
	if err != nil {
		Logger.Fatal(err)
	}

	Templates = templates
}
