// file: cmd/web/templates.go
package main

import (
	"html/template"
	"path/filepath"
)

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		fileName := filepath.Base(page)

		ts, err := template.ParseFiles("./ui/html/base.tmpl", page)
		if err != nil {
			return nil, err
		}

		cache[fileName] = ts
	}
	return cache, nil
}
