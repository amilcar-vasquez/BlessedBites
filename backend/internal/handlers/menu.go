package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/amilcar-vasquez/blessed-bites/backend/pkg/responses"
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
)

type MenuHandler struct {
	MenuItem *data.MenuItemModel
}

func NewMenuHandler(db *sql.DB) *MenuHandler {
	return &MenuHandler{MenuItem: &data.MenuItemModel{DB: db}}
}

func (h *MenuHandler) List(w http.ResponseWriter, r *http.Request) {
	if categoryRaw := r.URL.Query().Get("category"); categoryRaw != "" {
		categoryID, err := strconv.ParseInt(categoryRaw, 10, 64)
		if err != nil || categoryID <= 0 {
			responses.BadRequest(w, "invalid category filter")
			return
		}

		items, err := h.MenuItem.GetByCategoryID(categoryID)
		if err != nil {
			responses.InternalServerError(w, "failed to load menu")
			return
		}

		normalizeMenuImageURLs(items)

		responses.JSON(w, http.StatusOK, map[string]any{"items": items})
		return
	}

	activeOnly := r.URL.Query().Get("active") == "true"

	var (
		items []*data.MenuItem
		err   error
	)
	if activeOnly {
		items, err = h.MenuItem.GetAllActive()
	} else {
		items, err = h.MenuItem.GetAll()
	}
	if err != nil {
		responses.InternalServerError(w, "failed to load menu")
		return
	}

	normalizeMenuImageURLs(items)

	responses.JSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *MenuHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		responses.BadRequest(w, "invalid menu id")
		return
	}

	item, err := h.MenuItem.Get(id)
	if err != nil {
		if err == sql.ErrNoRows {
			responses.NotFound(w, "menu item not found")
			return
		}
		responses.InternalServerError(w, "failed to load menu item")
		return
	}

	item.ImageURL = normalizeImageURL(item.ImageURL)

	responses.JSON(w, http.StatusOK, item)
}

func (h *MenuHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		responses.BadRequest(w, "missing query parameter q")
		return
	}

	items, err := h.MenuItem.Search(q)
	if err != nil {
		responses.InternalServerError(w, "search failed")
		return
	}

	normalizeMenuImageURLs(items)

	responses.JSON(w, http.StatusOK, map[string]any{
		"query": q,
		"items": items,
	})
}

func normalizeMenuImageURLs(items []*data.MenuItem) {
	for _, item := range items {
		if item == nil {
			continue
		}
		item.ImageURL = normalizeImageURL(item.ImageURL)
	}
}

func normalizeImageURL(raw string) string {
	if raw == "" {
		return raw
	}

	normalized := strings.ReplaceAll(raw, "\\", "/")

	if idx := strings.Index(normalized, "/uploads/"); idx >= 0 {
		return normalized[idx:]
	}

	if idx := strings.Index(normalized, "uploads/"); idx >= 0 {
		return "/" + normalized[idx:]
	}

	if strings.HasPrefix(normalized, "http://") || strings.HasPrefix(normalized, "https://") || strings.HasPrefix(normalized, "/") {
		return normalized
	}

	return normalized
}
