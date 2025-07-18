// file: cmd/web/categoryHandlers.go
package main

import (
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/amilcar-vasquez/blessed-bites/internal/validator"
	"github.com/gorilla/csrf"
	"net/http"
	"strconv"
	"strings"
)

// POST handler to process the form submission for adding a new category
func (app *application) addCategory(w http.ResponseWriter, r *http.Request) {
	// Create a new template data instance

	app.logger.Info("Category form handler triggered")
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("Error parsing form data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Extract form fields
	name := r.PostForm.Get("category_name")

	// Create an instance of category
	category := &data.Category{
		Name: name,
	}

	// Validate the category data
	v := validator.NewValidator()
	data.ValidateCategory(v, category)
	if !v.ValidData() {
		app.logger.Error("Validation failed", "errors", v.Errors)
		data := app.addDefaultData(NewTemplateData(), w, r)
		data.Title = "Add Category"
		data.HeaderText = "Add Category"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"category_name": name,
		}
		data.CSRFField = csrf.TemplateField(r)
		err = app.render(w, http.StatusUnprocessableEntity, "AddMenuItem.tmpl", data)
		if err != nil {
			app.logger.Error("Error rendering template with validation errors", "error", err)
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

	categoryIDStr := r.FormValue("category_id")
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

func (app *application) viewMenuByCategory(w http.ResponseWriter, r *http.Request) {

	// Extract category ID from URL path
	idStr := strings.TrimPrefix(r.URL.Path, "/menu/category/")
	categoryID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Fetch category
	category, err := app.Category.GetByID(categoryID)
	if err != nil {
		app.logger.Error("Failed to get category", "error", err)
		http.NotFound(w, r)
		return
	}

	// Fetch all menu items for this category
	menuItems, err := app.MenuItem.GetByCategoryID(categoryID)
	if err != nil {
		app.logger.Error("Failed to get menu items by category", "error", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch all categories
	categories, err := app.Category.GetAll()
	if err != nil {
		app.logger.Error("Failed to get categories", "error", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	// Render the page
	data := app.addDefaultData(NewTemplateData(), w, r)
	data.MenuItems = menuItems
	data.Title = category.Name + " - Menu"
	data.Categories = categories

	err = app.render(w, http.StatusOK, "home.tmpl", data)
	if err != nil {
		app.logger.Error("Template rendering error", "error", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}
}
