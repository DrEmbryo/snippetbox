package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/DrEmbryo/snippetbox/cmd/pkg/forms"
	"github.com/DrEmbryo/snippetbox/cmd/pkg/models"
)

type templateData struct {
	AuthenticatedUser *models.User
	CSRFToken string
	CurrentYear int
	Form *forms.Form
	Flash string
	Snippet *models.Snippet
	Snippets []*models.Snippet
}

func humanDate(timestamp string) string {
	layout := "2006-01-02 15:04:05 -0700 MST"
	t, err := time.Parse(layout, timestamp)
	if err != nil {
		return timestamp
	}
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
	

func newTemplateCache(dir string) (map[string]*template.Template,  error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}