package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/gbih/snippetbox/pkg/forms"
	"github.com/gbih/snippetbox/pkg/models"
)

type templateData struct {
	CurrentYear int
	Flash       string
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        *forms.Form
	// FormData    url.Values        // access url.Values type
	// FormErrors  map[string]string // redisplay data upon errors
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// FuncMap is the type of the map defining the mapping from names to functions.
// There are 18 built-in template functions. To define a custom function:
// 1. Create a template.FuncMap object containing the custom humanDate() function.
// 2. Use the template.Funcs() method to register this before parsing the templates.

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method. This means we have to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		//ts, err := template.ParseFiles(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
