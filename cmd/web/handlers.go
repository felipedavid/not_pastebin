package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/felipedavid/not_pastebin/internal/models"
	"github.com/felipedavid/not_pastebin/internal/validator"
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

		flash := app.sessionManager.PopString(r.Context(), "flash")

		td := app.newTemplateData(r)
		td.Snippet = snippet
		td.Flash = flash

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

		app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	default:
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

type userSignupForm struct {
	Name                string `form:"name`
	Email               string `form:"email"`
	Password            string `form:"password`
	validator.Validator `form:"-`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := app.newTemplateData(r)
		data.Form = userSignupForm{}
		app.render(w, http.StatusOK, "signup.tmpl", data)
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		form := userSignupForm{
			Name:     r.PostForm.Get("name"),
			Email:    r.PostForm.Get("email"),
			Password: r.PostForm.Get("password"),
		}

		form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
		form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
		form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "This field must be a valid email address")
		form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
		form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

		if !form.Valid() {
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
			return
		}

		err = app.users.Insert(form.Name, form.Email, form.Password)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateEmail) {
				form.AddFieldError("email", "Email address is already in use")

				data := app.newTemplateData(r)
				data.Form = form
				app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
			} else {
				app.serverError(w, err)
			}

			return
		}

		app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	default:
		w.Header().Set("Allow", "GET, POST")
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "Form to create a new user")
	case http.MethodPost:
		fmt.Fprintf(w, "Creating a new user")
	default:
		w.Header().Set("Allow", "GET, POST")
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Fprintf(w, "Logging out from account")
	default:
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}
