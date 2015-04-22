package service

import (
	"errors"
	"log"
)

// MailerConfig ...
type MailerConfig struct {
	Type     string `xml:"type"`
	Host     string `xml:"host"`
	Username string `xml:"username"`
	Password string `xml:"password"`
	Port     int    `xml:"port"`
	From     string `xml:"from"`
}

const (
	// MailerTypeSMTP ...
	MailerTypeSMTP = "smtp"
)

// Mailer ...
type Mailer interface {
	SendWelcomeMail(to string, toUsername string) error
}

// Mail ...
var Mail Mailer

// InitMailer ...
func InitMailer(config MailerConfig) {
	if config.Type == MailerTypeSMTP {
		Mail = NewSMTP(config)
	} else {
		log.Fatalln(errors.New("Unsupported mailer type '" + config.Type + "'"))
	}
}
