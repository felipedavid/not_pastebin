package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/felipedavid/not_pastebin/internal/models"
)

func (a *app) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.notFound(w)
		return
	}

	snippets, err := a.snippets.Latest()
	if err != nil {
		a.serverError(w, err)
		return
	}

	for _, s := range snippets {
		fmt.Fprintf(w, "%+v", *s)
	}

	//files := []string{
	//	"./ui/html/base.tmpl",
	//	"./ui/html/pages/nav.tmpl",
	//	"./ui/html/pages/home.tmpl",
	//}

	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	a.serverError(w, err)
	//	return
	//}

	//err = ts.ExecuteTemplate(w, "base", nil)
	//if err != nil {
	//	a.serverError(w, err)
	//}
}

func (a *app) view(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			a.notFound(w)
			return
		}

		snippet, err := a.snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				a.notFound(w)
			} else {
				a.serverError(w, err)
			}
			return
		}

		data := &templateData{
			Snippet: snippet,
		}

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			"./ui/html/pages/view.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			a.serverError(w, err)
			return
		}

		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			a.serverError(w, err)
		}
	default:
		w.Header().Set("Allow", "GET")
		a.clientError(w, http.StatusMethodNotAllowed)
	}
}

func (a *app) create(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		_, err := a.snippets.Insert("Hello there", "No idea", 1)
		if err != nil {
			a.serverError(w, err)
		}
		return
	default:
		w.Header().Set("Allow", "POST")
		a.clientError(w, http.StatusMethodNotAllowed)
	}
}
