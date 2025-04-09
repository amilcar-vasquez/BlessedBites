package main

import (
	"net/http"
)

// BaseHandler renders the common elements
func (app *application) base(w http.ResponseWriter, r *http.Request) {
	// Prepare the common template data
	data := NewTemplateData()
	data.Title = "Blessed Bites - Base"
	data.HeaderText = "Welcome to Blessed Bites"

	// Render the base template
	err := app.render(w, http.StatusOK, "base.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render base template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Home page handler (now simplified)
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// Prepare the data for the home page (base template already has todos)
	data := NewTemplateData()
	data.Title = "Tapir Journals - Home"
	data.HeaderText = "Welcome to Tapir Journals"

	// Render the home template, which extends the base template
	err := app.render(w, http.StatusOK, "home.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
