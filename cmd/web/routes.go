package main

import (
	"net/http"
)

// TODO: Build some custom router or use a random package. Chaining handlers is getting kinda ugly.
func (a *app) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/", a.loadAndSave(a.authenticate(http.HandlerFunc(a.home))))
	mux.Handle("/about", a.loadAndSave(a.authenticate(http.HandlerFunc(a.about))))
	mux.Handle("/snippet/view/", a.loadAndSave(a.authenticate(http.HandlerFunc(a.view))))
	mux.Handle("/user/signup", a.loadAndSave(a.authenticate(http.HandlerFunc(a.signup))))
	mux.Handle("/user/login", a.loadAndSave(a.authenticate(http.HandlerFunc(a.login))))
	mux.Handle("/user/logout", a.loadAndSave(a.authenticate(http.HandlerFunc(a.logout))))

	mux.Handle("/snippet/create",
		a.sessionManager.LoadAndSave(a.authenticate(a.requireAuthentication(http.HandlerFunc(a.create)))))

	return a.recoverPanic(a.logRequest(secureHeaders(mux)))
}

func (a *app) loadAndSave(h http.Handler) http.Handler {
	return a.sessionManager.LoadAndSave(h)
}
