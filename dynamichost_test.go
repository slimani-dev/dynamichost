package dynamichost_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slimani-dev/dynamic-host"
)

func TestDynamicHost(t *testing.T) {
	cfg := dynamichost.CreateConfig()
	cfg.Headers = map[string]string{
		"Host": "^([^.]+)\\.localhost$ -> $1.example.com",
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
		t.Errorf("invalid host value: %s", req.Host)
	}
}
