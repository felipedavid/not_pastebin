package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorResponse(w, r, http.StatusNotFound)
		return
	}

	files := []string{
		"./ui/html/base.gohtml",
		"./ui/html/pages/home.gohtml",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError)
	}
}

func viewSnippetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getQueryInt(r, "id")
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Viewing the snippet %d", id)
}

func createSnippetHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Write([]byte("create a snippet"))
	default:
		methodNotAllowedResponse(w, r, http.MethodPost)
	}
}
