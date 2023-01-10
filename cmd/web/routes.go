package main

import (
	"net/http"
)

func (a *app) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/snippet/view", a.view)
	mux.HandleFunc("/snippet/create", a.create)

	return a.recoverPanic(a.logRequest(secureHeaders(mux)))
}
