package main

import (
	"net/http"
)

// TODO: Build some custom router or use a random package. Chaining handlers is getting kinda ugly.
func (a *app) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/", a.loadAndSave(a.home))
	mux.Handle("/about", a.loadAndSave(a.about))
	mux.Handle("/snippet/view/", a.loadAndSave(a.view))
	mux.Handle("/user/signup", a.loadAndSave(a.signup))
	mux.Handle("/user/login", a.loadAndSave(a.login))
	mux.Handle("/user/logout", a.loadAndSave(a.logout))

	mux.Handle("/snippet/create",
		a.sessionManager.LoadAndSave(a.requireAuthentication(http.HandlerFunc(a.create))))

	return a.recoverPanic(a.logRequest(secureHeaders(mux)))
}

func (a *app) loadAndSave(h func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return a.sessionManager.LoadAndSave(http.HandlerFunc(h))
}
