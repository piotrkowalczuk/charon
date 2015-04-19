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
	Logger LoggerConfig `xml:"logger"`
	DB     DBConfig     `xml:"database"`
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
