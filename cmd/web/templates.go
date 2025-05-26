// file: cmd/web/templates.go
package main

import (
	"encoding/json"
	"html/template"
	"path/filepath"
)

func (app *application) newTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"json": jsonFunc, // assumes jsonFunc is declared elsewhere
	}
}

func jsonFunc(v interface{}) (template.JS, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return template.JS(b), nil // use with care — safe if you're outputting to <script>
}

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
