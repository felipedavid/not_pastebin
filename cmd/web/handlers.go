package main

import (
	"errors"
	"github.com/felipedavid/not_pastebin/internal/models"
	"net/http"
	"strconv"
)

func (a *app) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		snippets, err := a.snippets.Latest()
		if err != nil {
			a.serverError(w, err)
		}

		a.render(w, http.StatusOK, "home.tmpl", &templateData{
			Snippets: snippets,
		})
	default:
		w.Header().Set("Allow", "GET")
		a.clientError(w, http.StatusMethodNotAllowed)
	}
}

func (a *app) snippet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		if err != nil {
			a.serverError(w, err)
			return
		}

		s, err := a.snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				a.notFound(w)
				return
			}
			a.serverError(w, err)
			return
		}

		a.render(w, http.StatusOK, "view.tmpl", &templateData{
			Snippet: s,
		})
	default:
		w.Header().Set("Allow", "GET")
		a.clientError(w, http.StatusMethodNotAllowed)
	}
}
