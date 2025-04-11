package main

import (
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"math/rand"
	"net/http"
	"time"
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

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//random menu item
	rand.Seed(time.Now().UnixNano())

	// Fetch categories for sidebar and buttons
	categories, err := app.Category.GetAll()
	if err != nil {
		app.logger.Error("failed to retrieve categories", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch menu items
	menuItems, err := app.MenuItem.GetAll()

	if err != nil {
		app.logger.Error("failed to retrieve menu items", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a random menu item
	var RandomMenuItems []*data.MenuItem
	if len(menuItems) > 0 {
		shuffled := make([]*data.MenuItem, len(menuItems))
		copy(shuffled, menuItems)
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})

		limit := 4
		if len(shuffled) < limit {
			limit = len(shuffled)
		}
		RandomMenuItems = shuffled[:limit]
	}

	// Prepare the data
	data := NewTemplateData()
	data.Title = "Welcome to Blessed Bites"
	data.HeaderText = "Welcome to Blessed Bites"
	data.Categories = categories
	data.MenuItems = menuItems
	data.RandomMenuItems = RandomMenuItems

	err = app.render(w, http.StatusOK, "home.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
