package main

import "net/http"

func (a *app) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/", a.loadAndSave(a.home))
	mux.Handle("/snippet", a.loadAndSave(a.snippet))
	mux.Handle("/snippet/create", a.loadAndSave(a.createSnippet))

	return a.recoverPanic(a.logRequest(a.secureHeaders(mux)))
}

func (a *app) loadAndSave(next func(http.ResponseWriter, *http.Request)) http.Handler {
	return a.sessionManager.LoadAndSave(http.HandlerFunc(next))
}
