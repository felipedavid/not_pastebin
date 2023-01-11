package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
		urlParameters := strings.Split(r.URL.Path, "/")[1:]
		fmt.Println(urlParameters)
		if len(urlParameters) < 3 {
			a.notFound(w)
			return
		}

		id, err := strconv.Atoi(urlParameters[2])
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
	case http.MethodGet:
		data := a.newTemplateData()
		a.render(w, http.StatusOK, "create.tmpl", data)
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
