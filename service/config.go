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
	Domain string `xml:"domain"`
	Server struct {
		Host     string `xml:"host"`
		Port     string `xml:"port"`
		TLS      bool   `xml:"tls"`
		CertFile string `xml:"cert-file"`
		KeyFile  string `xml:"key-file"`
	} `xml:"server"`
	Logger         LoggerConfig         `xml:"logger"`
	DB             DBConfig             `xml:"database"`
	Templates      TemplatesConfig      `xml:"templates"`
	Mailer         mailConfig           `xml:"email"`
	PasswordHasher PasswordHasherConfig `xml:"password"`
	Routing        RoutingConfig        `xml:"routing"`
	Translation    translationConfig    `xml:"translation"`
	Mnemosyne      MnemosyneConfig      `xml:"mnemosyne"`
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
