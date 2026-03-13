package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/amilcar-vasquez/blessed-bites/backend/internal/middleware"
	"github.com/amilcar-vasquez/blessed-bites/backend/pkg/responses"
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
)

type RatingsHandler struct {
	Ratings *data.RatingModel
}

func NewRatingsHandler(db *sql.DB) *RatingsHandler {
	return &RatingsHandler{Ratings: &data.RatingModel{DB: db}}
}

type CreateRatingRequest struct {
	MenuItemID int `json:"menu_item_id"`
	Rating     int `json:"rating"`
	UserID     int `json:"user_id"`
}

func (h *RatingsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateRatingRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequest(w, "invalid JSON payload")
		return
	}

	if !h.Ratings.ValidateRating(input.Rating) {
		responses.BadRequest(w, "rating must be between 1 and 5")
		return
	}

	if claims, ok := middleware.ClaimsFromContext(r.Context()); ok {
		input.UserID = int(claims.UserID)
	}
	if input.UserID <= 0 || input.MenuItemID <= 0 {
		responses.BadRequest(w, "menu_item_id and authenticated user are required")
		return
	}

	rating := &data.Rating{
		UserID:     input.UserID,
		MenuItemID: input.MenuItemID,
		Rating:     input.Rating,
	}

	if err := h.Ratings.Insert(rating); err != nil {
		responses.InternalServerError(w, "failed to submit rating")
		return
	}

	responses.JSON(w, http.StatusCreated, rating)
}

func (h *RatingsHandler) Average(w http.ResponseWriter, r *http.Request) {
	menuItemID, err := strconv.Atoi(r.PathValue("menu_item_id"))
	if err != nil || menuItemID <= 0 {
		responses.BadRequest(w, "invalid menu item id")
		return
	}

	avg, err := h.Ratings.GetAverageRating(menuItemID)
	if err != nil {
		responses.InternalServerError(w, "failed to load rating")
		return
	}

	responses.JSON(w, http.StatusOK, map[string]any{
		"menu_item_id": menuItemID,
		"average":      avg,
	})
}
