package service

import (
	"errors"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/go-soa/charon/lib/routing"
)

const (
	webDir  = "web"
	mailDir = "mail"
)

// ErrTemplateWrongDir ...
var ErrTemplateWrongDir = errors.New("service: can't use template from not existing directory")

// TemplatesConfig ...
type TemplatesConfig struct {
	Path string `xml:"path"`
}

// TemplateManager ...
type TemplateManager struct {
	templates map[string]*template.Template
}

// TplManager ...
var TplManager *TemplateManager

// InitTemplateManager ...
func InitTemplateManager(config TemplatesConfig) {

	tpls := map[string]*template.Template{}
	tpls[webDir] = initForDir(config.Path + "/" + webDir)
	tpls[mailDir] = initForDir(config.Path + "/" + mailDir)

	TplManager = &TemplateManager{
		templates: tpls,
	}
}

func initForDir(templateDir string) *template.Template {
	langDecider := func() (string, error) {
		return "en", nil
	}

	templates := template.New("")
	templates.Funcs(template.FuncMap{
		"url":     routing.URLTemplateFunc(URLGenerator),
		"url_abs": routing.URLAbsTemplateFunc(URLGenerator),
		"trans":   Translate.GetTransFunc(langDecider),
	})

	files := []string{}

	filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	templates, err := templates.ParseFiles(files...)
	if err != nil {
		log.Fatalln(err)
	}

	return templates
}

// Get ...
func (tm *TemplateManager) get(writer io.Writer, prefixDir string, templateName string, params map[string]interface{}) error {
	tpls, ok := tm.templates[prefixDir]
	if !ok {
		return ErrTemplateWrongDir
	}

	tpls.ExecuteTemplate(writer, templateName, params)

	return nil
}

// GetForMail ...
func (tm *TemplateManager) GetForMail(writer io.Writer, templateName string, params map[string]interface{}) error {
	return tm.get(writer, mailDir, templateName, params)
}

// GetForWeb ...
func (tm *TemplateManager) GetForWeb(writer io.Writer, templateName string, params map[string]interface{}) error {
	return tm.get(writer, webDir, templateName, params)
}
