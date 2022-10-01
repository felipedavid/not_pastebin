package main

import (
	"errors"
	"fmt"
	"github.com/felipedavid/not_pastebin/internal/models"
	"net/http"
	"strconv"
)

func (a *app) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.notFound(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		snippets, err := a.snippets.Latest()
		if err != nil {
			a.serverError(w, err)
		}

		data := newTemplateData(r)
		data.Snippets = snippets

		a.render(w, http.StatusOK, "home.tmpl", data)
	default:
		a.errorMethodNotAllowed(w, http.MethodGet)
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

		data := newTemplateData(r)
		data.Snippet = s

		a.render(w, http.StatusOK, "view.tmpl", data)
	default:
		a.errorMethodNotAllowed(w, http.MethodGet)
	}
}

func (a *app) createSnippet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := newTemplateData(r)
		a.render(w, http.StatusOK, "create_snippet.tmpl", data)
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			a.clientError(w, http.StatusBadRequest)
			return
		}

		title := r.PostForm.Get("title")
		content := r.PostForm.Get("content")

		expiresStr := r.PostForm.Get("expires")
		expires, err := strconv.ParseInt(expiresStr, 10, 64)
		if err != nil {
			a.clientError(w, http.StatusBadRequest)
			return
		}

		id, err := a.snippets.Insert(title, content, expires)
		if err != nil {
			a.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusMovedPermanently)
	default:
		a.errorMethodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}
