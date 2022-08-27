package main

import (
	"errors"
	"fmt"
	"github.com/felipedavid/not_pastebin/internal/models"
	"github.com/felipedavid/not_pastebin/internal/validator"
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

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := app.newTemplateData(r)
		data.Form = snippetCreateForm{
			Expires: 365,
		}
		app.render(w, http.StatusOK, "create.tmpl", data)
	case http.MethodPost:
		// Parse data from the request body and stores it
		// into r.PostForm map
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		expires, err := strconv.Atoi(r.PostForm.Get("expires"))
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		form := snippetCreateForm{
			Title:   r.PostForm.Get("title"),
			Content: r.PostForm.Get("content"),
			Expires: expires,
		}

		form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
		form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
		form.CheckField(validator.PermittedInt(form.Expires, 365, 7, 1), "expires", "This field must be equal to 1, 7 or 365")

		if !form.Valid() {
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
			return
		}

		id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
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
