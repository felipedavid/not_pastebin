package main

import (
	"errors"
	"fmt"
	"github.com/felipedavid/not_pastebin/internal/models"
	"github.com/felipedavid/not_pastebin/internal/validator"
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
			return
		}

		data := a.newTemplateData(r)
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

		data := a.newTemplateData(r)
		data.Snippet = s

		a.render(w, http.StatusOK, "view.tmpl", data)
	default:
		a.errorMethodNotAllowed(w, http.MethodGet)
	}
}

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (a *app) createSnippet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := a.newTemplateData(r)
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

		fieldErrors := make(map[string]string)
		val := validator.Validator{FieldErrors: fieldErrors}

		val.CheckField(validator.NotBlank(title), "title", "This field cannot be blank")
		val.CheckField(validator.MaxChars(title, 100), "title", "This field cannot be blank")
		val.CheckField(validator.NotBlank(content), content, "This field cannot be blank")
		val.CheckField(validator.PermittedInt(expires, 1, 7, 365), "expires",
			"This field must be equal to 1, 7 or 365")

		if !val.Valid() {
			data := a.newTemplateData(r)
			data.FieldErrors = fieldErrors
			a.render(w, http.StatusOK, "create_snippet.tmpl", data)
			return
		}

		id, err := a.snippets.Insert(title, content, expires)
		if err != nil {
			a.serverError(w, err)
			return
		}

		a.sessionManager.Put(r.Context(), "flash", "Snippet successfully created")
		http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusMovedPermanently)
	default:
		a.errorMethodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form::"-"`
}

func (a *app) userSignup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := a.newTemplateData(r)
		data.Form = userSignupForm{}
		a.render(w, http.StatusOK, "signup.tmpl", data)
	case http.MethodPost:
		var form userSignupForm

		err := a.decodePostForm(r, &form)
		if err != nil {
			a.clientError(w, http.StatusBadRequest)
			return
		}

		form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
		form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
		form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email")
		form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
		form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

		if !form.Valid() {
			data := a.newTemplateData(r)
			data.Form = form
			a.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
			return
		}

		err = a.users.Insert(form.Name, form.Email, form.Password)
		if err != nil {
			a.serverError(w, err)
		}

		a.sessionManager.Put(r.Context(), "flash", "Your signup was successul. Please log in.")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	default:
		a.errorMethodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func (a *app) userLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "This should be a html form for login")
	case http.MethodPost:
		fmt.Fprintf(w, "This should actually make the user login")
	default:
		a.errorMethodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func (a *app) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.errorMethodNotAllowed(w, http.MethodGet, http.MethodPost)
		return
	}

	fmt.Fprintf(w, "Logging out...")
}
