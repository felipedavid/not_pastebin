package main

import (
	"fmt"
	"github.com/felipedavid/not_pastebin/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// TODO: Automatically parse partials
		tFiles := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}

		ts, err := template.ParseFiles(tFiles...)
		if err != nil {
			return nil, err
		}

		filename := filepath.Base(page)
		cache[filename] = ts
	}

	return cache, nil
}

func (a *app) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := a.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		a.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		a.serverError(w, err)
	}
}
