package main

import "net/http"

// routes returns a serveMux with all our routes set up
func (a *app) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/snippet/view/", a.viewSnippet)
	mux.HandleFunc("/snippet/create", a.createSnippet)

	return mux
}
