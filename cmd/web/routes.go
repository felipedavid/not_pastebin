package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	sm := app.sessionManager

	mux.Handle("/", sm.LoadAndSave(http.HandlerFunc(app.home)))
	mux.Handle("/snippet/view/", sm.LoadAndSave(http.HandlerFunc(app.snippetView)))
	mux.Handle("/snippet/create", sm.LoadAndSave(http.HandlerFunc(app.snippetCreate)))

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
