package main

import (
	"fmt"
	"net/http"
)

func (a *app) logRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		re := newResponseWriter(w)
		next.ServeHTTP(re, r)

		a.infoLogger.Printf("%s \"%s %s\" -> %d %s\n", r.RemoteAddr, r.Method, r.URL.Path, re.status, http.StatusText(re.status))
	}
	return http.HandlerFunc(fn)
}

// secureHeaders setups the secure headers recommend by OWASP
func (a *app) secureHeaders(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (a *app) recoverPanic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				a.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)

}
