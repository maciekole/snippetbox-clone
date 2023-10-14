package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.maciekole.net/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
}

// Create a humanDate function which return a nicely formatted string repr of a time.Time object
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Init new cache
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/html/pages/*.tmpl". This will essentially gives
	// us a slice of all the filepaths for our application 'page' templates
	// like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// register template.FuncMap
		templateSet, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		templateSet, err = templateSet.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the  page template.
		templateSet, err = templateSet.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map as normal...
		cache[name] = templateSet
	}

	return cache, nil
}
