package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/amilcar-vasquez/blessed-bites/backend/pkg/responses"
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
)

type AdminHandler struct {
	DB       *sql.DB
	MenuItem *data.MenuItemModel
	Category *data.CategoryModel
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{
		DB:       db,
		MenuItem: &data.MenuItemModel{DB: db},
		Category: &data.CategoryModel{DB: db},
	}
}

type UpsertMenuItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  int     `json:"category_id"`
	ImageURL    string  `json:"image_url"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type AdminOrder struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	FullName  string    `json:"full_name"`
	TotalCost float64   `json:"total_cost"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *AdminHandler) ListMenu(w http.ResponseWriter, r *http.Request) {
	items, err := h.MenuItem.GetAll()
	if err != nil {
		responses.InternalServerError(w, "failed to load menu")
		return
	}
	responses.JSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *AdminHandler) CreateMenu(w http.ResponseWriter, r *http.Request) {
	var input UpsertMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequest(w, "invalid JSON payload")
		return
	}

	item := &data.MenuItem{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CategoryID:  input.CategoryID,
		ImageURL:    input.ImageURL,
	}
	if err := h.MenuItem.Insert(item); err != nil {
		responses.BadRequest(w, "failed to create menu item")
		return
	}

	if input.IsActive != nil {
		_ = h.MenuItem.SetActiveState(item.ID, *input.IsActive)
	}

	responses.JSON(w, http.StatusCreated, item)
}

func (h *AdminHandler) UpdateMenu(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		responses.BadRequest(w, "invalid menu id")
		return
	}

	existing, err := h.MenuItem.Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responses.NotFound(w, "menu item not found")
			return
		}
		responses.InternalServerError(w, "failed to load menu item")
		return
	}

	var input UpsertMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequest(w, "invalid JSON payload")
		return
	}

	existing.Name = input.Name
	existing.Description = input.Description
	existing.Price = input.Price
	existing.CategoryID = input.CategoryID
	existing.ImageURL = input.ImageURL

	if err := h.MenuItem.Update(existing); err != nil {
		responses.InternalServerError(w, "failed to update menu item")
		return
	}

	if input.IsActive != nil {
		if err := h.MenuItem.SetActiveState(existing.ID, *input.IsActive); err != nil {
			responses.InternalServerError(w, "failed to update menu active state")
			return
		}
	}

	responses.JSON(w, http.StatusOK, existing)
}

func (h *AdminHandler) DeleteMenu(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		responses.BadRequest(w, "invalid menu id")
		return
	}

	if err := h.MenuItem.Delete(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responses.NotFound(w, "menu item not found")
			return
		}
		responses.InternalServerError(w, "failed to delete menu item")
		return
	}

	responses.JSON(w, http.StatusOK, map[string]any{"deleted": id})
}

func (h *AdminHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var input CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequest(w, "invalid JSON payload")
		return
	}

	category := &data.Category{Name: input.Name}
	if err := h.Category.Insert(category); err != nil {
		responses.BadRequest(w, "failed to create category")
		return
	}

	responses.JSON(w, http.StatusCreated, category)
}

func (h *AdminHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		responses.BadRequest(w, "invalid category id")
		return
	}

	if err := h.Category.Delete(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responses.NotFound(w, "category not found")
			return
		}
		responses.InternalServerError(w, "failed to delete category")
		return
	}

	responses.JSON(w, http.StatusOK, map[string]any{"deleted": id})
}

func (h *AdminHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT o.id, o.user_id, u.full_name, o.total_cost, o.status, o.created_at
		FROM orders o
		JOIN users u ON u.id = o.user_id
		ORDER BY o.created_at DESC
		LIMIT 200
	`

	rows, err := h.DB.Query(query)
	if err != nil {
		responses.InternalServerError(w, "failed to load orders")
		return
	}
	defer rows.Close()

	items := make([]AdminOrder, 0)
	for rows.Next() {
		var item AdminOrder
		if err := rows.Scan(&item.ID, &item.UserID, &item.FullName, &item.TotalCost, &item.Status, &item.CreatedAt); err != nil {
			responses.InternalServerError(w, "failed to parse orders")
			return
		}
		items = append(items, item)
	}

	responses.JSON(w, http.StatusOK, map[string]any{"items": items})
}
