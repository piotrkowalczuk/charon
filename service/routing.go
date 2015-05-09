package service

import (
	"encoding/xml"

	"github.com/go-soa/charon/lib/routing"
)

// Routes ...
var Routes routing.Routes

// URLGenerator ...
var URLGenerator routing.URLGenerator

type RoutingConfig struct {
	IsSSL       bool      `xml:"is-ssl"`
	BaseURL     string    `xml:"base-url"`
	RoutesGroup XMLRoutes `xml:"routes"`
}

type XMLRoutes struct {
	XMLName xml.Name   `xml:"routes"`
	Routes  []XMLRoute `xml:"route"`
}

type XMLRoute struct {
	Name    string `xml:"name"`
	Pattern string `xml:"pattern"`
}

func (rc *RoutingConfig) RoutesMap() map[string]string {
	result := make(map[string]string)

	for _, route := range rc.RoutesGroup.Routes {
		result[route.Name] = route.Pattern
	}

	return result
}

// InitRouting ...
func InitRouting(config RoutingConfig) {
	Routes = routing.NewRoutes(config.RoutesMap())
	URLGenerator = routing.NewURLGenerator(config.BaseURL, config.IsSSL, Routes)
}
