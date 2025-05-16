package main

import (
	"encoding/json"
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"net/http"
	"strconv"
	"time"
)

// create raiting handler
func (app *application) createRatingHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID     int `json:"user_id"`
		MenuItemID int `json:"menu_item_id"`
		Rating     int `json:"rating"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.logger.Error("Failed to decode JSON", "error", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	rating := &data.Rating{
		UserID:     input.UserID,
		MenuItemID: input.MenuItemID,
		Rating:     input.Rating,
		CreatedAt:  time.Now(),
	}

	// Validate the rating
	if !app.Rating.ValidateRating(rating.Rating) {
		http.Error(w, "Invalid rating. Rating must be between 1 and 5.", http.StatusBadRequest)
		return
	}

	// Insert the rating into the database
	err = app.Rating.Insert(rating)
	if err != nil {
		app.logger.Error("Failed to insert rating", "error", err)
		http.Error(w, "Failed to insert rating", http.StatusInternalServerError)
		return
	}
	// Return the created rating as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rating)

}

// get rating handler
func (app *application) getAverageRatingHandler(w http.ResponseWriter, r *http.Request) {
	menuItemIDStr := r.URL.Query().Get("menu_item_id")
	if menuItemIDStr == "" {
		http.Error(w, "Missing menu_item_id parameter", http.StatusBadRequest)
		return
	}

	menuItemID, err := strconv.Atoi(menuItemIDStr)
	if err != nil {
		http.Error(w, "Invalid menu_item_id parameter. It must be an integer.", http.StatusBadRequest)
		return
	}

	averageRating, err := app.Rating.GetAverageRating(menuItemID)
	if err != nil {
		app.logger.Error("Failed to get average rating", "error", err)
		http.Error(w, "Failed to get average rating", http.StatusInternalServerError)
		return
	}

	response := struct {
		MenuItemID    string  `json:"menu_item_id"`
		AverageRating float64 `json:"average_rating"`
	}{
		MenuItemID:    strconv.Itoa(menuItemID),
		AverageRating: averageRating,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
