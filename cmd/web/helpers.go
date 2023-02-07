package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
)

func (a *app) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.errLogger.Output(2, trace)
	if a.debugMode {
		http.Error(w, trace, http.StatusInternalServerError)
		return
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *app) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (a *app) notFound(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}

// render gets the template from the template cache, executes it, and write the result to the client
func (a *app) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// if the application is in debug mode, don't use the template cache
	var temp *template.Template
	if a.debugMode {
		ts, err := template.New(page).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			a.serverError(w, err)
			return
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			a.serverError(w, err)
			return
		}

		ts, err = ts.ParseFiles(fmt.Sprintf("./ui/html/pages/%s", page))
		if err != nil {
			a.serverError(w, err)
			return
		}
		temp = ts
	} else {
		ts, ok := a.templateCache[page]
		if !ok {
			err := fmt.Errorf("the template %s does not exist", page)
			a.serverError(w, err)
			return
		}
		temp = ts
	}

	buf := new(bytes.Buffer)

	err := temp.ExecuteTemplate(buf, "base", data)
	if err != nil {
		a.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (a *app) isAuthenticated(r *http.Request) bool {
	authenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return authenticated
}
