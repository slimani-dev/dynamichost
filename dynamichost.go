package dynamichost

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"text/template"
)

// Config holds the plugin configuration.
type Config struct {
	Headers []HeaderConfig `json:"headers,omitempty"`
}

// HeaderConfig represents a single header transformation rule.
type HeaderConfig struct {
	Name         string `json:"name"`
	RegexPattern string `json:"regexPattern"`
	NewHost      string `json:"newHost"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Headers: []HeaderConfig{},
	}
}

// DynamicHost is the plugin structure.
type DynamicHost struct {
	next     http.Handler
	headers  []HeaderConfig
	name     string
	template *template.Template
}

// New creates a new DynamicHost plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Headers) == 0 {
		return nil, fmt.Errorf("headers configuration cannot be empty")
	}

	return &DynamicHost{
		headers:  config.Headers,
		next:     next,
		name:     name,
		template: template.New("dynamichost").Delims("[[", "]]"),
	}, nil
}

// ServeHTTP processes the request and modifies the Host header accordingly.
func (dh *DynamicHost) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, header := range dh.headers {
		if header.Name == "Host" {
			match, _ := regexp.MatchString(header.RegexPattern, req.Host)
			if match {
				tmpl, err := dh.template.Parse(header.NewHost)
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}

				writer := &bytes.Buffer{}
				err = tmpl.Execute(writer, req)
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}

				req.Host = writer.String()
				req.Header.Set("Host", writer.String())
			}
		}
	}
	dh.next.ServeHTTP(rw, req)
}
