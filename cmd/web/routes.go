package main

import (
	"net/http"
)

func (a *app) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/", a.loadAndSave(a.home))
    mux.Handle("/about", a.loadAndSave(a.about))
	mux.Handle("/snippet/view/", a.loadAndSave(a.view))
	mux.Handle("/snippet/create", a.loadAndSave(a.create))

	return a.recoverPanic(a.logRequest(secureHeaders(mux)))
}

func (a *app) loadAndSave(h func(w http.ResponseWriter, r *http.Request)) http.Handler {
    return a.sessionManager.LoadAndSave(http.HandlerFunc(h))
}
