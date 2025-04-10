// filepath: cmd/web/main.go
package main

import (
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/amilcar-vasquez/blessed-bites/internal/validator"
	"net/http"
	"strconv"
)

// POST handler to process the form submission for adding a new category

func (app *application) addCategory(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("Error parsing form data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Extract form fields
	name := r.PostForm.Get("name")

	// Create an instance of category
	category := &data.Category{
		Name: name,
	}

	// Validate the category data
	v := validator.NewValidator()
	data.ValidateCategory(v, category)
	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Add Category"
		data.HeaderText = "Add Category"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"name": name,
		}
		err := app.render(w, http.StatusUnprocessableEntity, "AddCategory.tmpl", data)
		if err != nil {
			app.logger.Error("Error rendering template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}
	// Insert the category into the database
	err = app.Category.Insert(category)
	if err != nil {
		app.logger.Error("Error inserting category", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Redirect to the menu page
	http.Redirect(w, r, "/menu/add", http.StatusSeeOther)
}

// delete handler
func (app *application) deleteCategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	categoryIDStr := r.FormValue("id")
	if categoryIDStr == "" {
		app.logger.Error("missing category ID", "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		app.logger.Error("invalid category ID format", "category_id", categoryIDStr, "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := app.Category.Delete(categoryID); err != nil {
		app.logger.Error("failed to delete category entry", "error", err, "category_id", categoryID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/menu/add", http.StatusSeeOther)
}
