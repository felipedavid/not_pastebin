package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFoundResponse(w)
		return
	}

	files := []string{
		"./ui/html/base.gohtml",
		"./ui/html/pages/home.gohtml",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app *application) viewSnippetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getQueryInt(r, "id")
	if err != nil {
		app.badRequestResponse(w)
		return
	}
	fmt.Fprintf(w, "Viewing the snippet %d", id)
}

func (app *application) createSnippetHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Write([]byte("create a snippet"))
	default:
		methodNotAllowedResponse(w, r, http.MethodPost)
	}
}
