package main

import "net/http"

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hMap := w.Header()
		hMap.Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		hMap.Set("Referrer-Policy", "origin-when-cross-origin")
		hMap.Set("X-Content-Type-Options", "nosniff")
		hMap.Set("X-Frame-Options", "deny")
		hMap.Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (a *app) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.infoLogger.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}
