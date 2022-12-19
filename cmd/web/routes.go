package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vancanhuit/snippetbox/ui"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		app.clientError(w, http.StatusMethodNotAllowed)
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	r.Use(app.recoverPanic)
	r.Use(app.logRequest)
	r.Use(secureHeaders)

	fileServer := http.FileServer(http.FS(ui.Files))
	r.Handle("/static/*", fileServer)

	r.Get("/ping", ping)

	dynamic := []func(next http.Handler) http.Handler{
		app.sessionManager.LoadAndSave,
		noSurf,
		app.authenticate,
	}
	protected := append(dynamic, app.requireAuthentication)

	r.With(dynamic...).Get("/", app.home)
	r.With(dynamic...).Get("/about", app.about)
	r.With(dynamic...).Get(
		"/snippet/view/{id}", app.snippetView)
	r.With(protected...).Get(
		"/snippet/create", app.snippetCreateView)
	r.With(protected...).Post(
		"/snippet/create", app.snippetCreatePost)

	r.With(dynamic...).Get(
		"/user/signup", app.userSignupView)
	r.With(dynamic...).Post(
		"/user/signup", app.userSignupPost)
	r.With(dynamic...).Get(
		"/user/login", app.userLoginView)
	r.With(dynamic...).Post(
		"/user/login", app.userLoginPost)
	r.With(protected...).Post(
		"/user/logout", app.userLogoutPost)
	r.With(protected...).Get("/account/view", app.accountView)
	r.With(protected...).Get(
		"/account/password/update", app.accountPasswordUpdateView)
	r.With(protected...).Post(
		"/account/password/update", app.accountPasswordUpdatePost)

	return r
}
