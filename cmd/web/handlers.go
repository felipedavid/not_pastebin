package main

import (
	"errors"
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

	data := a.newTemplateData()
	data.Snippets = snippets

	a.render(w, http.StatusOK, "home.tmpl", data)
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

		data := a.newTemplateData()
		data.Snippet = snippet

		a.render(w, http.StatusOK, "view.tmpl", data)
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
