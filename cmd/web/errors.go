package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverErrorResponse(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func (app *application) clientErrorResponse(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFoundResponse(w http.ResponseWriter) {
	app.clientErrorResponse(w, http.StatusNotFound)
}

func (app *application) badRequestResponse(w http.ResponseWriter) {
	app.clientErrorResponse(w, http.StatusBadRequest)
}
