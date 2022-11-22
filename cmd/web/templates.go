package main

import (
	"html/template"
	"path/filepath"

	"github.com/felipedavid/not_pastebin/internal/models"
)

// templateData will hold any dynamic data that we pass to our templates
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
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

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			filePath,
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[fileName] = ts
	}

	return cache, err
}
