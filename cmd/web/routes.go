package main

import "net/http"

func (a *app) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/snippet", a.snippet)

	return a.logRequest(mux)
}
