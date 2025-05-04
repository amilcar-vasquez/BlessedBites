// file: cmd/web/menuHandlers.go
package main

import (
	"fmt"
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/amilcar-vasquez/blessed-bites/internal/validator"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func generateFileName(categoryID int, name, originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	if ext == "" {
		ext = ".jpg"
	}

	cleanName := strings.ToLower(name)
	cleanName = strings.ReplaceAll(cleanName, " ", "-")
	cleanName = strings.ReplaceAll(cleanName, "_", "-")
	cleanName = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, cleanName)

	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("%s_%s_%s%s", timestamp, strconv.Itoa(categoryID), cleanName, ext)
}

// parseMenuItemForm parses the multipart form and returns a populated MenuItem and form values.
func (app *application) parseMenuItemForm(r *http.Request, isMultipart bool) (*data.MenuItem, map[string]string, map[string]string, error) {
	var formErrors = make(map[string]string)
	var formData = make(map[string]string)

	// Parse form
	if isMultipart {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			return nil, formData, formErrors, fmt.Errorf("parsing multipart form: %w", err)
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			return nil, formData, formErrors, fmt.Errorf("parsing form: %w", err)
		}
	}

	// Extract fields
	idStr := r.FormValue("id") // only for updates
	name := r.FormValue("name")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	categoryIDStr := r.FormValue("category_id")

	// Save raw form values
	formData["id"] = idStr
	formData["name"] = name
	formData["description"] = description
	formData["price"] = priceStr
	formData["category_id"] = categoryIDStr

	// Handle price
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		formErrors["price"] = "Invalid price format"
	}

	// Handle category ID
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		formErrors["category_id"] = "Invalid category ID"
	}

	// Handle file upload if multipart
	var imageURL string
	if isMultipart {
		file, header, err := r.FormFile("image")
		if err != nil && err != http.ErrMissingFile {
			return nil, formData, formErrors, fmt.Errorf("retrieving uploaded file: %w", err)
		}
		if file != nil {
			defer file.Close()
			fileName := generateFileName(categoryID, name, header.Filename)
			imagePath := "./ui/static/img/uploads/" + fileName
			dst, err := os.Create(imagePath)
			if err != nil {
				return nil, formData, formErrors, fmt.Errorf("saving uploaded file: %w", err)
			}
			defer dst.Close()
			if _, err := io.Copy(dst, file); err != nil {
				return nil, formData, formErrors, fmt.Errorf("copying uploaded file: %w", err)
			}
			imageURL = imagePath
		}
	}

	// Build the menu item
	menuItem := &data.MenuItem{
		Name:        name,
		Description: description,
		Price:       price,
		CategoryID:  categoryID,
		ImageURL:    imageURL,
	}

	return menuItem, formData, formErrors, nil
}

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
	menuItem, formData, formErrors, err := app.parseMenuItemForm(r, true)
	if err != nil {
		app.logger.Error("Error parsing form", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	v := validator.NewValidator()
	data.ValidateMenuItem(v, menuItem)
	for k, vErr := range v.Errors {
		formErrors[k] = vErr
	}

	if len(formErrors) > 0 {
		data := NewTemplateData()
		data.Title = "Add Menu Item"
		data.HeaderText = "Add Menu Item"
		data.FormErrors = formErrors
		data.FormData = formData
		categories, err := app.Category.GetAll()
		if err != nil {
			app.logger.Error("Error retrieving categories from database", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		data.Categories = categories

		err = app.render(w, http.StatusUnprocessableEntity, "AddMenuItem.tmpl", data)
		if err != nil {
			app.logger.Error("Error rendering template with errors", "error", err)
		}
		return
	}

	err = app.MenuItem.Insert(menuItem)
	if err != nil {
		app.logger.Error("Error inserting menu item", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

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

	menuItemIDStr := r.FormValue("id")
	if menuItemIDStr == "" {
		app.logger.Error("missing ID", "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	MenuItemID, err := strconv.ParseInt(menuItemIDStr, 10, 64)
	if err != nil {
		app.logger.Error("invalid ID format", "id", menuItemIDStr, "error", err)
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
	if err := r.ParseForm(); err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	idStr := r.FormValue("id")
	// Check if ID is present in the form data

	if idStr == "" {
		app.logger.Error("Menu item ID is missing in form data")
		http.Error(w, "Missing menu item ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		app.logger.Error("Missing menu item ID", "value", idStr, "error", err)
		http.Error(w, "Where is", http.StatusBadRequest)
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
	menuItem, formData, formErrors, err := app.parseMenuItemForm(r, true)
	if err != nil {
		app.logger.Error("Error parsing form", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get ID separately
	idStr := formData["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.logger.Error("Invalid menu item ID", "value", idStr, "error", err)
		http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
		return
	}
	menuItem.ID = int64(id)

	v := validator.NewValidator()
	data.ValidateMenuItem(v, menuItem)
	for k, vErr := range v.Errors {
		formErrors[k] = vErr
	}

	if len(formErrors) > 0 {
		data := NewTemplateData()
		data.Title = "Edit Menu Item"
		data.HeaderText = "Edit Menu Item"
		data.FormErrors = formErrors
		data.FormData = formData
		data.MenuItem = menuItem

		err := app.render(w, http.StatusUnprocessableEntity, "AddMenuItem.tmpl", data)
		if err != nil {
			app.logger.Error("Error rendering template with errors", "error", err)
		}
		return
	}

	err = app.MenuItem.Update(menuItem)
	if err != nil {
		app.logger.Error("Error updating menu item", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/menu", http.StatusSeeOther)
}
