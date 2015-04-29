package routing

import (
	"fmt"
	"net/url"
	"strings"
)

type URLGenerator struct {
	baseURL string
	routes  Routes
}

// NewURLGenerator ...
func NewURLGenerator(baseURL string, routes Routes) URLGenerator {
	return URLGenerator{
		baseURL: baseURL,
		routes:  routes,
	}
}

func (ug *URLGenerator) generate(path string, params map[string]interface{}) (*url.URL, error) {
	for paramName, paramValue := range params {
		strings.Replace(path, "{"+paramName+"}", fmt.Sprintf("%s", paramValue), -1)
	}

	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Generate generates relative path. Naive implementation.
func (ug *URLGenerator) Generate(routeName RouteName, params map[string]interface{}) (*url.URL, error) {
	path := ug.routes.GetPattern(routeName).String()

	return ug.generate(path, params)
}

// GenerateAbs generates absolute path. Naive implementation.
func (ug *URLGenerator) GenerateAbs(routeName RouteName, params map[string]interface{}) (*url.URL, error) {
	path := ug.routes.GetPattern(routeName).String()

	return ug.generate(ug.baseURL+path, params)
}
