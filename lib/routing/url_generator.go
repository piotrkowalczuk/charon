package routing

import (
	"fmt"
	"net/url"
	"strings"
)

type URLGenerator struct {
	isSSL   bool
	baseURL string
	routes  Routes
}

// NewURLGenerator ...
func NewURLGenerator(baseURL string, isSSL bool, routes Routes) URLGenerator {
	return URLGenerator{
		isSSL:   isSSL,
		baseURL: baseURL,
		routes:  routes,
	}
}

func (ug *URLGenerator) generate(path string, params map[string]interface{}) (*url.URL, error) {
	for paramName, paramValue := range params {
		// TODO: add more types if necessary
		switch p := paramValue.(type) {
		case int:
			path = strings.Replace(path, ":"+paramName, fmt.Sprintf("%d", p), -1)
		case int64:
			path = strings.Replace(path, ":"+paramName, fmt.Sprintf("%d", p), -1)
		case string:
			path = strings.Replace(path, ":"+paramName, p, -1)
		}
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
	prefix := "http://"
	if ug.isSSL {
		prefix = "https://"
	}

	path := ug.routes.GetPattern(routeName).String()

	return ug.generate(prefix+ug.baseURL+path, params)
}
