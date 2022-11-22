package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/felipedavid/not_pastebin/internal/models"
)

// home is a handler. Handlers in go are like controllers in the MVC pattern
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

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		a.serverError(w, err)
		return
	}

	data := &templateData{Snippets: snippets}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		a.serverError(w, err)
	}

}

func (a *app) viewSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		a.notFound(w)
		return
	}

	s, err := a.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			a.notFound(w)
		} else {
			a.serverError(w, err)
		}
		return
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

	data := &templateData{Snippet: s}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) createSnippet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		title := "Test"
		content := "This is supposed to be just a test snippet"
		expires := 7

		id, err := a.snippets.Insert(title, content, expires)
		if err != nil {
			a.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	default:
		a.clientError(w, http.StatusMethodNotAllowed)
	}
}
