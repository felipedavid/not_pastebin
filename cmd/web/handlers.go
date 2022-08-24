package main

import (
	"errors"
	"fmt"
	"github.com/felipedavid/not_pastebin/internal/models"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path != "/" {
			app.notFound(w)
			return
		}

		snippets, err := app.snippets.Latest()
		if err != nil {
			app.serverError(w, err)
			return
		}

		td := app.newTemplateData(r)
		td.Snippets = snippets

		app.render(w, http.StatusOK, "home.tmpl", td)
	default:
		w.Header().Set("Allow", "GET")
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id, err := strconv.Atoi(getParameter(r.URL.Path, 2))
		if err != nil || id < 1 {
			app.notFound(w)
			return
		}

		snippet, err := app.snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}

		td := app.newTemplateData(r)
		td.Snippet = snippet

		app.render(w, http.StatusOK, "view.tmpl", td)
	default:
		w.Header().Set("Allow", "GET")
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "Display the form for creating new snippets")
	case http.MethodPost:
		title := "O snail"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
		expires := 7

		id, err := app.snippets.Insert(title, content, expires)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	default:
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}
