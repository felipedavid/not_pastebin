package main

import (
	"net/http"
	"runtime/debug"
	"strings"
)

func (a *app) serverError(w http.ResponseWriter, err error) {
	a.errLogger.Printf("%s\n%s", err.Error(), debug.Stack())
	a.render(w, http.StatusInternalServerError, "error.tmpl", &templateData{
		ErrorCode: http.StatusInternalServerError,
	})
}

func (a *app) clientError(w http.ResponseWriter, code int) {
	a.render(w, code, "error.tmpl", &templateData{ErrorCode: code})
}

func (a *app) notFound(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}

func (a *app) errorMethodNotAllowed(w http.ResponseWriter, allowed ...string) {
	w.Header().Set("Allow", strings.Join(allowed, ", "))
	a.clientError(w, http.StatusMethodNotAllowed)
}

func (a *app) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = a.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		return err
	}

	return nil
}
