package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/", app.loadAndSaveMiddleware(app.home))
	mux.Handle("/snippet/view/", app.loadAndSaveMiddleware(app.snippetView))
	mux.Handle("/snippet/create", app.loadAndSaveMiddleware(app.snippetCreate))
	mux.Handle("/user/signup", app.loadAndSaveMiddleware(app.userSignup))
	mux.Handle("/user/login", app.loadAndSaveMiddleware(app.userLogin))
	mux.Handle("/user/logout", app.loadAndSaveMiddleware(app.userLogout))

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}

func (app *application) loadAndSaveMiddleware(f http.HandlerFunc) http.Handler {
	return app.sessionManager.LoadAndSave(http.HandlerFunc(f))
}
