package main

import "net/http"

// responseWriter is just a wrapper around the http's ResponseWriter to enable
// us to log the status code of the response after writing to it
type responseWriter struct {
	http.ResponseWriter
	status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (re *responseWriter) WriteHeader(status int) {
	re.status = status
	re.ResponseWriter.WriteHeader(status)
}
