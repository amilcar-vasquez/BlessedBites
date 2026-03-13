package handlers

import (
	"database/sql"
	"net/http"

	"github.com/amilcar-vasquez/blessed-bites/backend/pkg/responses"
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
)

type CategoriesHandler struct {
	Category *data.CategoryModel
}

func NewCategoriesHandler(db *sql.DB) *CategoriesHandler {
	return &CategoriesHandler{Category: &data.CategoryModel{DB: db}}
}

func (h *CategoriesHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.Category.GetAll()
	if err != nil {
		responses.InternalServerError(w, "failed to load categories")
		return
	}
	responses.JSON(w, http.StatusOK, map[string]any{"items": categories})
}
