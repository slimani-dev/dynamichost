// Package dynamichost provides a middleware plugin for Traefik that dynamically rewrites
// request headers based on configurable regex patterns.
package dynamichost

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"text/template"
)

// HeaderConfig defines the structure for header transformations.
type HeaderConfig struct {
	Name         string `json:"name"`
	RegexPattern string `json:"regexPattern"`
	NewHost      string `json:"newHost"`
}

// Config the plugin configuration.
type Config struct {
	Headers []HeaderConfig `json:"headers,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Headers: []HeaderConfig{},
	}
}

// DynamicHost a plugin that rewrites headers dynamically.
type DynamicHost struct {
	next     http.Handler
	headers  []HeaderConfig
	name     string
	template *template.Template
}

// New creates a new DynamicHost plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Headers) == 0 {
		return nil, fmt.Errorf("headers cannot be empty")
	}

	return &DynamicHost{
		headers:  config.Headers,
		next:     next,
		name:     name,
		template: template.New("dynamichost").Delims("[[", "]]"),
	}, nil
}

func (a *DynamicHost) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, header := range a.headers {
		re, err := regexp.Compile(header.RegexPattern)
		if err != nil {
			http.Error(rw, "Invalid regex pattern", http.StatusInternalServerError)
			return
		}

		newHost := re.ReplaceAllString(req.Host, header.NewHost)
		req.Host = newHost
		req.Header.Set("Host", newHost)
	}

	a.next.ServeHTTP(rw, req)
}
