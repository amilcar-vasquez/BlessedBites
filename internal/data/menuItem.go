// filepath: internal/data/menuItem.go
package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/amilcar-vasquez/blessed-bites/internal/validator"
)

// MenuItem struct to hold the data for a menu item
type MenuItem struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CategoryID  int       `json:"category_id"`
	OrderCount  int       `json:"order_count"`
	IsActive    bool      `json:"is_active"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}

// MenuModel struct to hold the database connection pool
type MenuItemModel struct {
	DB *sql.DB
}

// Insert inserts a new menu item into the database
func (m *MenuItemModel) Insert(menuItem *MenuItem) error {
	query := `INSERT INTO menu_items (name, description, price, category_id, image_url) 
			  VALUES ($1, $2, $3, $4, $5) 
			  RETURNING id, order_count, is_active, created_at`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(
		ctx,
		query,
		menuItem.Name,
		menuItem.Description,
		menuItem.Price,
		menuItem.CategoryID,
		menuItem.ImageURL,
	).Scan(&menuItem.ID, &menuItem.OrderCount, &menuItem.IsActive, &menuItem.CreatedAt)
}

// ValidateMenuItem validates the fields of a menu item
func ValidateMenuItem(v *validator.Validator, item *MenuItem) {
	v.Check(validator.NotBlank(item.Name), "name", "Name is required")
	v.Check(validator.MaxLength(item.Name, 100), "name", "Name must be less than 100 characters")

	v.Check(item.Price > 0, "price", "Price must be greater than zero")

	v.Check(validator.MaxLength(item.Description, 500), "description", "Description must be less than 500 characters")
	v.Check(item.CategoryID >= 0, "category_id", "Category ID must be a non-negative integer")
}

// GetAll retrieves all menu items from the database
func (m *MenuItemModel) GetAll() ([]*MenuItem, error) {
	query := `SELECT id, name, description, price, category_id, order_count, is_active, image_url, created_at 
			  FROM menu_items 
			  ORDER BY created_at DESC`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*MenuItem{}
	for rows.Next() {
		item := &MenuItem{}
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.CategoryID, &item.OrderCount, &item.IsActive, &item.ImageURL, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// Delete deletes a menu item from the database
func (m *MenuItemModel) Delete(id int64) error {
	query := `DELETE FROM menu_items WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
