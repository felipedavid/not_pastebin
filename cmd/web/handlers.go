package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/felipedavid/not_pastebin/internal/models"
	"github.com/felipedavid/not_pastebin/internal/validator"
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

	data := a.newTemplateData(r)
	data.Snippets = snippets

	a.render(w, http.StatusOK, "home.tmpl", data)
}

func (a *app) view(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		urlParameters := strings.Split(r.URL.Path, "/")[1:]
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

		data := a.newTemplateData(r)
		data.Snippet = snippet

		a.render(w, http.StatusOK, "view.tmpl", data)
	default:
		w.Header().Set("Allow", "GET")
		a.clientError(w, http.StatusMethodNotAllowed)
	}
}

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
    validator.Validator
}

func (a *app) create(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := a.newTemplateData(r)
		a.render(w, http.StatusOK, "create.tmpl", data)
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			a.clientError(w, http.StatusBadRequest)
			return
		}

		expires, err := strconv.Atoi(r.PostForm.Get("expires"))
		if err != nil {
			a.clientError(w, http.StatusBadRequest)
			return
		}

		form := snippetCreateForm{
			Title:       r.PostForm.Get("title"),
			Content:     r.PostForm.Get("content"),
			Expires:     expires,
		}

        form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
        form.CheckField(validator.MaxChars(form.Title, 100), "title", 
            "This field cannot be more than 100 characters long")
        form.CheckField(validator.NotBlank(form.Content), "content", 
            "This field cannot be blank")
        form.CheckField(validator.PermittedInt(expires, 1, 7, 365), "expires", 
            "This field must be equal to 1, 7 or 365")

		if !form.Valid() {
			data := a.newTemplateData(r)
			data.Form = form
			a.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
			return
		}

		id, err := a.snippets.Insert(form.Title, form.Content, form.Expires)
		if err != nil {
			a.serverError(w, err)
			return
		}

        a.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	default:
		w.Header().Set("Allow", "POST")
		a.clientError(w, http.StatusMethodNotAllowed)
	}
}
