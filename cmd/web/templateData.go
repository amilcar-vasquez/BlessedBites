// file: cmd/web/templateData.go
package main

import (
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
)

type TemplateData struct {
	Title           string
	HeaderText      string
	MenuItems       []*data.MenuItem
	MenuItem        *data.MenuItem
	Categories      []*data.Category
	CategoryMap     map[int]string
	RandomMenuItems []*data.MenuItem

	FormErrors map[string]string
	FormData   map[string]string
}

// factory function to initialize a new templateData struct
func NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:      "Welcome to Blessed Bites",
		HeaderText: "Welcome to Blessed Bites",
		FormErrors: map[string]string{},
		FormData:   map[string]string{},
	}
}
