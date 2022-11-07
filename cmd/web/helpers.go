package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError writes an error message and stack trace to the errLogger
// then sends a 500 error message to the client
func (a *app) serverError(w http.ResponseWriter, err error) {
	// We are using debug.Stack() to get the stack trace of the current goroutine
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.errLogger.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends a error message with specific statusCode to the client
func (a *app) clientError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

// notFound sends a 404 Not Found response to the client
func (a *app) notFound(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}
