package service

import (
	"encoding/xml"
	"log"
	"os"
)

// Config ...
var Config AppConfig

// AppConfig ...
type AppConfig struct {
	Server struct {
		Host string `xml:"host"`
		Port string `xml:"port"`
	} `xml:"server"`
	Logger         LoggerConfig         `xml:"logger"`
	DB             DBConfig             `xml:"database"`
	Template       TemplateConfig       `xml:"template"`
	Mailer         MailerConfig         `xml:"email"`
	PasswordHasher PasswordHasherConfig `xml:"password"`
}

// InitConfig ...
func InitConfig(environment string) {
	file, err := os.Open("conf/" + environment + ".xml")
	if err != nil {
		log.Fatalf("Cannot open configuration file: %v\n", err)
	}

	defer file.Close()
	decoder := xml.NewDecoder(file)

	decoder.Decode(&Config)
}
