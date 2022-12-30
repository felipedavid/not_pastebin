package main

import (
	"fmt"
	"net/http"
	"html/template"
)

func (a *app) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
        a.notFound(w)
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
        a.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
        a.serverError(w, err)
	}
}

func (a *app) view(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Viewing a specific snippet")
}

func (a *app) create(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Fprintf(w, "Creating a snippet")
		return
	default:
		w.Header().Set("Allow", "POST")
		w.WriteHeader(405)
		a.clientError(w, http.StatusMethodNotAllowed)
	}
}

