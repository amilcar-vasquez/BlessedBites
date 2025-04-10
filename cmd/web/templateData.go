// file: cmd/web/templateData.go
package main

import (
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
)

type TemplateData struct {
	Title      string
	HeaderText string
	MenuItems  []*data.MenuItem
	Categories []*data.Category

	FormErrors map[string]string
	FormData   map[string]string
}

// factory function to initialize a new templateData struct
func NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:      "Welcome to Tapir Journals",
		HeaderText: "Welcome to Tapir Journals",
		FormErrors: map[string]string{},
		FormData:   map[string]string{},
	}
}
