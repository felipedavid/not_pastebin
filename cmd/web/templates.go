package main

import (
	"bytes"
	"fmt"
	"github.com/felipedavid/not_pastebin/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	ErrorCode   int
	FieldErrors map[string]string
	Flash       string
}

func newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

func statusCodeText(status int) string {
	return http.StatusText(status)
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate":      humanDate,
	"statusCodeText": statusCodeText,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		filename := filepath.Base(page)

		ts, err := template.New(filename).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
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

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		a.serverError(w, err)
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}
