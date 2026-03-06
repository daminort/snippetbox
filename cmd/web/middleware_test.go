package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.demien.net/internal/assert"
)

func TestCommonHeaders(t *testing.T) {
	app := &Application{
		logger: slog.New(slog.DiscardHandler),
	}

	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(app.ping)
	commonHeaders(next).ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	headers := map[string]string{
		"Content-Security-Policy": "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		"Referrer-Policy":         "origin-when-cross-origin",
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "deny",
		"X-XSS-Protection":        "0",
		"X-Server":                "Go",
	}

	for k, v := range headers {
		assert.Equal(t, res.Header.Get(k), v)
	}

	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")

}
