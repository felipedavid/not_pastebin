package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.homeHandler)
	mux.HandleFunc("/snippet/view", app.viewSnippetHandler)
	mux.HandleFunc("/snippet/create", app.createSnippetHandler)

	return mux
}
