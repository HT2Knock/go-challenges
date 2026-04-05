package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test array for convenience
var testCases = []struct {
	name       string
	method     string
	url        string
	token      string
	wantStatus int
	wantBody   string
}{
	{
		name:       "Public /healthz endpoint with no token",
		method:     "GET",
		url:        "/healthz",
		token:      "",
		wantStatus: http.StatusOK,
		wantBody:   http.StatusText(http.StatusOK),
	},
	{
		name:       "Secure /secure endpoint no token",
		method:     "GET",
		url:        "/secure",
		token:      "",
		wantStatus: http.StatusUnauthorized,
		wantBody:   http.StatusText(http.StatusUnauthorized),
	},
	{
		name:       "Secure /secure endpoint invalid token",
		method:     "GET",
		url:        "/secure",
		token:      "invalid",
		wantStatus: http.StatusUnauthorized,
		wantBody:   http.StatusText(http.StatusUnauthorized),
	},
	{
		name:       "Secure /secure endpoint correct token",
		method:     "GET",
		url:        "/secure",
		token:      "secret",
		wantStatus: http.StatusOK,
		wantBody:   http.StatusText(http.StatusOK),
	},
	{
		name:       "Public /healthz endpoint with invalid token",
		method:     "GET",
		url:        "/healthz",
		token:      "wrong",
		wantStatus: http.StatusOK,
		wantBody:   http.StatusText(http.StatusOK),
	},
	{
		name:       "Public /healthz endpoint with correct token",
		method:     "GET",
		url:        "/healthz",
		token:      "secret",
		wantStatus: http.StatusOK,
		wantBody:   http.StatusText(http.StatusOK),
	},
	{
		name:       "Different method on /secure with valid token",
		method:     "POST",
		url:        "/secure",
		token:      "secret",
		wantStatus: http.StatusOK,
		wantBody:   http.StatusText(http.StatusOK),
	},
	{
		name:       "Different method on /secure with no token",
		method:     "POST",
		url:        "/secure",
		token:      "",
		wantStatus: http.StatusUnauthorized,
		wantBody:   http.StatusText(http.StatusUnauthorized),
	},
}

func TestMiddleware(t *testing.T) {
	server := NewServer(slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, nil)
			if tc.token != "" {
				req.Header.Set("X-Auth-Token", tc.token)
			}
			rr := httptest.NewRecorder()

			server.ServeHTTP(rr, req)

			if rr.Code != tc.wantStatus {
				t.Errorf("Expected status %d, got %d", tc.wantStatus, rr.Code)
			}
			body := strings.TrimSpace(rr.Body.String())
			if body != tc.wantBody {
				t.Errorf("Expected body %q, got %q", tc.wantBody, body)
			}
		})
	}
}

func BenchmarkSecureRoute(b *testing.B) {
	server := NewServer(slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)))

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/secure", nil)
		req.Header.Set("X-Auth-Token", "secret")
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)
		io.Copy(io.Discard, rr.Result().Body)
	}
}
