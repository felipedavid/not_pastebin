package main

import (
	"net/http"
	"runtime/debug"
)

func (a *app) serverError(w http.ResponseWriter, err error) {
	a.errLogger.Printf("%s\n%s", err.Error(), debug.Stack())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *app) clientError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
