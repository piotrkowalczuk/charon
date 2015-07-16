package service

import (
	"log"

	"github.com/go-soa/charon/lib"
)

const (
	registrationConfirmationMailerDirectory = "registration_confirmation"
	passwordRecoveryMailerDirectory         = "password_recovery"
)

var (
	// RegistrationConfirmationMailer ...
	RegistrationConfirmationMailer lib.Sender
	// PasswordRecoveryMailer ...
	PasswordRecoveryMailer lib.Sender
)

type mailConfig struct {
	Type     string `xml:"type"`
	Host     string `xml:"host"`
	Username string `xml:"username"`
	Password string `xml:"password"`
	Port     int    `xml:"port"`
	From     string `xml:"from"`
}

// InitMailers ...
func InitMailers(config mailConfig, tplMgr *TemplateManager) {
	var transport lib.Transporter
	switch config.Type {
	case lib.TransporterTypeSMTP:
		transport = lib.NewSMTPTransporter(config.Host, config.Username, config.Password, config.Port)
	default:
		log.Fatalf("Unsupported mailer type '%s'", config.Type)
	}

	registrationConfirmationMailer, err := lib.NewMailer(registrationConfirmationMailerDirectory, config.From, transport, tplMgr)
	if err != nil {
		Logger.Fatal(err)
	}

	passwordRecoveryMailer, err := lib.NewMailer(passwordRecoveryMailerDirectory, config.From, transport, tplMgr)
	if err != nil {
		Logger.Fatal(err)
	}

	RegistrationConfirmationMailer = registrationConfirmationMailer
	PasswordRecoveryMailer = passwordRecoveryMailer
}
