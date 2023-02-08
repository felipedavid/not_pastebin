package main

import (
	"bytes"
	"github.com/felipedavid/not_pastebin/internal/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	resRec := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(resRec, req)

	res := resRec.Result()

	// Check if all the headers were set
	assert.Equal(t, res.Header.Get("Content-Security-Policy"),
		"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
	assert.Equal(t, res.Header.Get("Referrer-Policy"), "origin-when-cross-origin")
	assert.Equal(t, res.Header.Get("X-Content-Type-Options"), "nosniff")
	assert.Equal(t, res.Header.Get("X-Frame-Options"), "deny")
	assert.Equal(t, res.Header.Get("X-XSS-Protection"), "0")

	assert.Equal(t, res.StatusCode, http.StatusOK)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
