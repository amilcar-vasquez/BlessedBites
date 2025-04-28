// file: cmd/web/menuHandlers.go
package main

import (
	"fmt"
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/amilcar-vasquez/blessed-bites/internal/validator"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// GET handler to display the form to add a new menu item
func (app *application) addMenuItemForm(w http.ResponseWriter, r *http.Request) {

	// Create a new template data instance
	data := NewTemplateData()
	data.Title = "Add Menu Item"
	data.HeaderText = "Add Menu Item"

	// Retrieve all categories from the database
	categories, err := app.Category.GetAll()
	if err != nil {
		app.logger.Error("Error retrieving categories from database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data.Categories = categories

	// Render the add menu item form template
	err = app.render(w, http.StatusOK, "AddMenuItem.tmpl", data)
	if err != nil {
		app.logger.Error("Error rendering template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// POST handler to process the form submission for adding a new menu item
func (app *application) addMenuItemHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		app.logger.Error("Error parsing multipart form data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Extract form fields
	name := r.FormValue("name")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	categoryIDStr := r.FormValue("category_id")

	// handle file upload
	file, header, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		app.logger.Error("Error retrieving file from form", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if file != nil {
		defer file.Close()
	}
	// Save the file to a specific location
	var imageURL string
	if file != nil {
		fileName := fmt.Sprintf("%s.%s.%s.%s.%s.%s_%s_%s_%s",
			time.Now().Format("2006"),
			time.Now().Format("01"),
			time.Now().Format("02"),
			time.Now().Format("15"),
			time.Now().Format("04"),
			time.Now().Format("05"),
			categoryIDStr,
			name,
			header.Filename)
		imageURL = "./ui/static/img/uploads/" + fileName
		dst, err := os.Create(imageURL)
		if err != nil {
			app.logger.Error("Error creating file", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, file)
		if err != nil {
			app.logger.Error("Error saving file", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	}

	// Convert price and category ID to appropriate types
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		app.logger.Error("Error converting price to float", "error", err)
		http.Error(w, "Invalid Price", http.StatusBadRequest)
		return
	}
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		app.logger.Error("Error converting category ID to int", "error", err)
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// Create an instance of MenuItem
	menuItem := &data.MenuItem{
		Name:        name,
		Description: description,
		Price:       price,
		CategoryID:  categoryID,
		ImageURL:    imageURL,
	}

	// Validate the menu item data
	v := validator.NewValidator()
	data.ValidateMenuItem(v, menuItem)
	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Add Menu Item"
		data.HeaderText = "Add Menu Item"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"name":        name,
			"description": description,
			"price":       priceStr,
			"category_id": categoryIDStr,
		}
		// Re-render the form with validation errors
		err := app.render(w, http.StatusUnprocessableEntity, "AddMenuItem.tmpl", data)
		if err != nil {
			app.logger.Error("Error rendering template with errors", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	// Insert the menu item into the database
	err = app.MenuItem.Insert(menuItem)
	if err != nil {
		app.logger.Error("Error inserting menu item into database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Redirect to the menu items list page
	http.Redirect(w, r, "/menu", http.StatusSeeOther)

}

// GET handler to display the list of menu items
func (app *application) menuPageHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve all menu items from the database
	menuItems, err := app.MenuItem.GetAll()
	if err != nil {
		app.logger.Error("Error retrieving menu items from database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//also retrieve all categories
	categories, err := app.Category.GetAll()

	//create a category map
	categoryMap := make(map[int]string)
	for _, category := range categories {
		categoryMap[int(category.ID)] = category.Name
	}

	if err != nil {
		app.logger.Error("Error retrieving categories from database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a new template data instance
	data := NewTemplateData()
	data.Title = "Menu Items"
	data.HeaderText = "Menu Items"
	data.MenuItems = menuItems
	data.Categories = categories
	data.CategoryMap = categoryMap

	// Render the menu items list template
	err = app.render(w, http.StatusOK, "menu.tmpl", data)
	if err != nil {
		app.logger.Error("Error rendering template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// delete handler
func (app *application) deleteMenuItem(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	MenuItemIDStr := r.FormValue("id")
	if MenuItemIDStr == "" {
		app.logger.Error("missing ID", "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	MenuItemID, err := strconv.ParseInt(MenuItemIDStr, 10, 64)
	if err != nil {
		app.logger.Error("invalid ID format", "id", MenuItemIDStr, "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := app.MenuItem.Delete(MenuItemID); err != nil {
		app.logger.Error("failed to delete Menu entry", "error", err, "id", MenuItemID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/menu", http.StatusSeeOther)
}

// GET handler to display the form to edit a menu item
func (app *application) editMenuItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)

	if idStr == "" {
		app.logger.Error("Menu item ID is missing in form data")
		http.Error(w, "Missing menu item ID", http.StatusBadRequest)
		return
	}

	if err != nil {
		app.logger.Error("Invalid menu item ID", "value", idStr, "error", err)
		http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
		return
	}

	menuItem, err := app.MenuItem.Get(int64(id))
	if err != nil {
		app.logger.Error("Menu item not found", "error", err)
		http.Error(w, "Menu Item Not Found", http.StatusNotFound)
		return
	}

	categories, err := app.Category.GetAll()
	if err != nil {
		app.logger.Error("Error retrieving categories", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData()
	data.Title = "Edit Menu Item"
	data.HeaderText = "Edit Menu Item"
	data.MenuItem = menuItem
	data.Categories = categories

	err = app.render(w, http.StatusOK, "AddMenuItem.tmpl", data)
	if err != nil {
		app.logger.Error("Error rendering template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// POST handler to process the form submission for editing a menu item
func (app *application) updateMenuItem(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("Error parsing form data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Extract form fields
	idStr := r.PostForm.Get("id")
	name := r.PostForm.Get("name")
	description := r.PostForm.Get("description")
	priceStr := r.PostForm.Get("price")

	categoryIDStr := r.PostForm.Get("category_id")
	app.logger.Info("Received price string", "price", priceStr)

	// Convert price and category ID to appropriate types
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		app.logger.Error("Error converting price to float", "error", err)
		http.Error(w, "Invalid Price", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		app.logger.Error("Error converting category ID to int", "error", err)
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.logger.Error("Error converting menu item ID to int", "error", err)
		http.Error(w, "Invalid Menu Item ID", http.StatusBadRequest)
		return
	}

	// Create an instance of MenuItem
	menuItem := &data.MenuItem{
		ID:          int64(id),
		Name:        name,
		Description: description,
		Price:       price,
		CategoryID:  categoryID,
	}

	v := validator.NewValidator()
	data.ValidateMenuItem(v, menuItem)
	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Update Menu Item"
		data.HeaderText = "Update Menu Item"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"name":        name,
			"description": description,
			"price":       priceStr,
			"category_id": categoryIDStr,
			"id":          idStr,
		}

		err = app.render(w, http.StatusUnprocessableEntity, "AddMenuItem.tmpl", data)
		if err != nil {
			app.logger.Error("Error rendering template with errors", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}
	// Update the menu item in the database
	err = app.MenuItem.Update(menuItem)
	if err != nil {
		app.logger.Error("Error updating menu item in database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Redirect to the menu items list page
	http.Redirect(w, r, "/menu", http.StatusSeeOther)
}
