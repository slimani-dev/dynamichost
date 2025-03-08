package dynamichost_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slimani-dev/dynamichost"
)

func TestDynamicHost(t *testing.T) {
	cfg := dynamichost.CreateConfig()
	cfg.Headers = []dynamichost.HeaderConfig{
		{
			Name:         "Host",
			RegexPattern: "^([^.]+)\\.localhost$",
			NewHost:      "$1.example.com",
		},
	}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := dynamichost.New(ctx, next, cfg, "dynamic-host-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://sub.localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeader(t, req, "Host", "sub.example.com")
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Host != expected {
		t.Errorf("expected %s to be %s, but got %s", key, expected, req.Host)
	}
}
