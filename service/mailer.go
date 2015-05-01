package service

import (
	"log"

	"github.com/go-soa/charon/mail"
)

// ConfirmationMailer ...
var ConfirmationMailer mail.Sender

type mailConfig struct {
	Type     string `xml:"type"`
	Host     string `xml:"host"`
	Username string `xml:"username"`
	Password string `xml:"password"`
	Port     int    `xml:"port"`
	From     string `xml:"from"`
}

// InitMailer ...
func InitMailer(config mailConfig, tplMgr *TemplateManager) {
	var transport mail.Transporter
	switch config.Type {
	case mail.TransporterTypeSMTP:
		transport = mail.NewSMTPTransporter(config.Host, config.Username, config.Password, config.Port)
	default:
		log.Fatalf("Unsupported mailer type '%s'", config.Type)
	}

	confirmationMailer, err := mail.NewConfirmationMailer(config.From, transport, tplMgr)
	if err != nil {
		Logger.Fatal(err)
	}

	ConfirmationMailer = confirmationMailer
}
