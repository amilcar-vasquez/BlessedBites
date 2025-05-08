// filepath: internal/data/menuItem.go
package data

import (
	"context"
	"database/sql"
	"strings"
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
	v.Check(validator.MinLength(item.Name, 3), "name", "Name must be at least 3 characters")
	v.Check(validator.MaxLength(item.Name, 100), "name", "Name must be less than 50 characters")

	v.Check(item.Price > 0, "price", "Price must be greater than zero")

	v.Check(validator.NotBlank(item.Description), "description", "Description is required")
	v.Check(validator.MaxLength(item.Description, 500), "description", "Description must be less than 500 characters")
	v.Check(validator.MinLength(item.Description, 10), "description", "Description must be at least 10 characters")

	v.Check(item.CategoryID >= 0, "category_id", "Please select a category")
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

// update updates a menu item in the database
func (m *MenuItemModel) Update(item *MenuItem) error {
	query := `UPDATE menu_items 
			  SET name = $1, description = $2, price = $3, category_id = $4, image_url = $5 
			  WHERE id = $6`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(
		ctx,
		query,
		item.Name,
		item.Description,
		item.Price,
		item.CategoryID,
		item.ImageURL,
		item.ID,
	)
	return err
}

// Get retrieves a menu item by ID from the database
func (m *MenuItemModel) Get(id int64) (*MenuItem, error) {
	query := `SELECT id, name, description, price, category_id, order_count, is_active, image_url, created_at 
			  FROM menu_items 
			  WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	item := &MenuItem{}
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.CategoryID, &item.OrderCount, &item.IsActive, &item.ImageURL, &item.CreatedAt)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (m *MenuItemModel) Search(query string) ([]*MenuItem, error) {
	query = "%" + strings.ToLower(query) + "%"

	stmt := `
		SELECT id, name, description, price, category_id, image_url
		FROM menu_items
		WHERE LOWER(name) LIKE $1 OR LOWER(description) LIKE $2
		ORDER BY name
	`

	rows, err := m.DB.Query(stmt, query, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*MenuItem
	for rows.Next() {
		var item MenuItem
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.CategoryID, &item.ImageURL)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (m *MenuItemModel) GetByCategoryID(categoryID int64) ([]*MenuItem, error) {
	query := `
		SELECT id, name, description, price, image_url, category_id
		FROM menu_items
		WHERE category_id = $1
		ORDER BY name ASC
	`

	rows, err := m.DB.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*MenuItem
	for rows.Next() {
		var item MenuItem
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.ImageURL, &item.CategoryID)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}
