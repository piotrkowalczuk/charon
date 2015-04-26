package service

import (
	"errors"
	"html/template"
	"log"

	"github.com/go-soa/auth/mail"
)

// Mail ...
var Mail *mail.Mail

// InitMailer ...
func InitMailer(config mailConfig, templates *template.Template) {
	var transport mail.Transporter
	if config.Type == mailTransportSMTP {
		transport = mail.NewTransportSMTP(config.Host, config.Username, config.Password, config.Port)
	} else {
		log.Fatalln(errors.New("Unsupported mailer type '" + config.Type + "'"))
		return
	}

	Mail = mail.NewMail(transport, config.From, templates)
}

type mailConfig struct {
	Type     string `xml:"type"`
	Host     string `xml:"host"`
	Username string `xml:"username"`
	Password string `xml:"password"`
	Port     int    `xml:"port"`
	From     string `xml:"from"`
}

const (
	// TypeSMTP ...
	mailTransportSMTP = "smtp"
)
