package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/felipedavid/not_pastebin/internal/models"
)

type templateData struct {
	CurrentYear       int
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
	Form              any
	Flash             string
	AuthenticatedUser bool
	User              *models.User
}

func (a *app) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		Flash:       a.sessionManager.PopString(r.Context(), "flash"),
		Form: snippetCreateForm{
			Expires: 1,
		},
		AuthenticatedUser: a.isAuthenticated(r),
	}
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// newTemplateCache returns a map containing all templates pre-parsed
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
