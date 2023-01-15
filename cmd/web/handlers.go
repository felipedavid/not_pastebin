package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

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
		err := r.ParseForm()
		if err != nil {
			a.clientError(w, http.StatusBadRequest)
			return
		}

		fieldErrors := make(map[string]string)

		title := r.PostForm.Get("title")
		content := r.PostForm.Get("content")
		expires, err := strconv.Atoi(r.PostForm.Get("expires"))
		if err != nil {
			a.clientError(w, http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(title) == "" {
			fieldErrors["title"] = "This field cannot be blank"
		} else if utf8.RuneCountInString(title) > 100 {
			fieldErrors["title"] = "This field cannot be more than 100 characters long"
		}

		if strings.TrimSpace(content) == "" {
			fieldErrors["content"] = "This field cannot be blank"
		}

		if expires != 1 && expires != 7 && expires != 365 {
			fieldErrors["expires"] = "This field must be equal to 1, 7 or 365"
		}

		if len(fieldErrors) > 0 {
			fmt.Fprint(w, fieldErrors)
			return
		}

		id, err := a.snippets.Insert(title, content, expires)
		if err != nil {
			a.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	default:
		w.Header().Set("Allow", "POST")
		a.clientError(w, http.StatusMethodNotAllowed)
	}
}
