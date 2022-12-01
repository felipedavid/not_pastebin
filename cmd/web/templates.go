package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/felipedavid/not_pastebin/internal/models"
)

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// templateData will hold any dynamic data that we pass to our templates
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func (a *app) newTemplateData(r *http.Request) *templateData {
	return &templateData{CurrentYear: time.Now().Year()}
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Get a slice of filepathes for all the files that match
	// the regular expression
	files, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, filePath := range files {
		fileName := filepath.Base(filePath)

		ts, err := template.New(fileName).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(filePath)
		if err != nil {
			return nil, err
		}

		cache[fileName] = ts
	}

	return cache, err
}
