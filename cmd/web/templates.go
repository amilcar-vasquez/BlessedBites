// file: cmd/web/templates.go
package main

import (
	"encoding/json"
	"github.com/amilcar-vasquez/blessed-bites/internal/utils"
	"html/template"
	"path/filepath"
)

// create a template func map to hold custom functions
var templateFuncMap = template.FuncMap{
	"json":  jsonFunc,
	"add":   utils.Add,
	"sub":   utils.Subtract,
	"until": utils.Until,
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

		ts, err := template.New("base.tmpl").Funcs(templateFuncMap).ParseFiles("./ui/html/base.tmpl", page)
		if err != nil {
			return nil, err
		}

		cache[fileName] = ts
	}
	return cache, nil
}
