package main

import (
	"errors"
	"fmt"
	"github.com/felipedavid/not_pastebin/internal/data"
	"net/http"
	"strconv"
)

func (a *app) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		snippets, err := a.snippets.Latest()
		if err != nil {
			a.serverError(w, err)
		}

		for _, s := range snippets {
			fmt.Fprintf(w, "%v\n", *s)
		}
	default:
		w.Header().Set("Allow", "GET")
		a.clientError(w, http.StatusMethodNotAllowed)
	}
	//tFiles := []string{
	//	"./ui/html/base.tmpl",
	//	"./ui/html/partials/nav.tmpl",
	//}

	//ts, err := template.ParseFiles(tFiles...)
	//if err != nil {
	//	a.serverError(w, err)
	//	return
	//}

	//err = ts.ExecuteTemplate(w, "base", nil)
	//if err != nil {
	//	a.serverError(w, err)
	//	return
	//}
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
			if errors.Is(err, data.ErrNoRecord) {
				a.notFound(w)
				return
			}
			a.serverError(w, err)
			return
		}

		fmt.Fprintf(w, "%v", *s)
	default:
		w.Header().Set("Allow", "GET")
		a.clientError(w, http.StatusMethodNotAllowed)
	}
}
