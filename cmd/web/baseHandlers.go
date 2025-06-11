package main

import (
	"encoding/json"
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/amilcar-vasquez/blessed-bites/internal/utils"
	"net/http"
	"sort"
	"strconv"
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

// helper function to get top user items:
func (app *application) getTopUserMenuItems(userID int, limit int) []*data.MenuItem {
	topItems := []*data.MenuItem{}
	topIDs, err := app.Recommendation.GetTopRecommendationsByUser(userID, limit)
	if err != nil {
		app.logger.Error("Failed to get user recommendations", "error", err)
		return topItems
	}
	for _, id := range topIDs {
		item, err := app.MenuItem.Get(int64(id))
		if err != nil {
			app.logger.Error("Failed to fetch menu item", "error", err, "id", id)
			continue
		}
		topItems = append(topItems, item)
	}
	return topItems
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	app.logger.Info("Initial page value", "page", page)

	// Fetch categories for sidebar and buttons
	categories, err := app.Category.GetAll()
	if err != nil {
		app.logger.Error("failed to retrieve categories", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch menu items
	menuItems, err := app.MenuItem.GetAllActive()

	if err != nil {
		app.logger.Error("failed to retrieve menu items", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Inject popularity flags BEFORE rendering
	popularIDs, err := app.Recommendation.GetPopularItemIDs(3)
	if err != nil {
		app.logger.Error("Error retrieving popular items", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	popularSet := make(map[int]struct{})
	for _, id := range popularIDs {
		popularSet[id] = struct{}{}
	}
	for i := range menuItems {
		if _, found := popularSet[int(menuItems[i].ID)]; found {
			menuItems[i].Popular = true
		}
	}

	// Reorder menuItems so that popular items appear first
	sort.SliceStable(menuItems, func(i, j int) bool {
		if menuItems[i].Popular && !menuItems[j].Popular {
			return true // i comes before j
		}
		return false
	})

	pagedMenuItems, currentPage, totalPages := utils.Paginate(menuItems, page, 9) // Use the Paginate function from utils

	// Output log with values for debugging
	app.logger.Info("Pagination info",
		"requestedPage", page,
		"currentPage", currentPage,
		"totalPages", totalPages,
		"totalMenuItems", len(menuItems),
		"pagedMenuItemsCount", len(pagedMenuItems),
	) // Use the Paginate function from utils

	//call the getTopUserMenuItems function to get the top user menu items
	var TopUserMenuItems []*data.MenuItem
	user := app.contextGetUser(r)
	if user == nil {
		app.logger.Info("User not found in session")
	} else {
		TopUserMenuItems = app.getTopUserMenuItems(int(user.ID), 3)
	}

	// Prepare the data
	data := app.addDefaultData(NewTemplateData(), w, r) // <- THIS LINE IS KEY
	data.Title = "Welcome to Blessed Bites"
	data.HeaderText = "Welcome to Blessed Bites"
	data.Categories = categories
	data.MenuItems = pagedMenuItems
	data.CurrentPage = currentPage
	data.TotalPages = totalPages
	data.TopUserMenuItems = TopUserMenuItems

	err = app.render(w, http.StatusOK, "home.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) searchMenuJSONHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	menuItems, err := app.MenuItem.Search(query)
	if err != nil {
		app.logger.Error("failed to search menu items", "error", err)
		http.Error(w, "Failed to search menu items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(menuItems)
	if err != nil {
		app.logger.Error("failed to encode JSON response", "error", err)
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	session, err := app.sessionStore.Get(r, "session")
	if err != nil {
		return nil
	}

	authenticated, ok := session.Values["authenticated"].(bool)
	if !ok || !authenticated {
		return nil
	}

	userID, ok := session.Values["authenticatedUserID"].(int64)
	if !ok || userID == 0 {
		return nil
	}

	user, err := app.User.GetByID(userID)
	if err != nil {
		return nil
	}

	return user
}
