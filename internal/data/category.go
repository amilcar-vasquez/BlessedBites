// data/categories.go
package data

import (
	"context"
	"database/sql"
	"github.com/amilcar-vasquez/blessed-bites/internal/validator"
	"time"
)

// Category struct to hold the data for a category
type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// CategoryModel struct to hold the database connection pool
type CategoryModel struct {
	DB *sql.DB
}

// Insert inserts a new category into the database
func (c *CategoryModel) Insert(category *Category) error {
	query := `INSERT INTO categories (name) 
			  VALUES ($1) 
			  RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return c.DB.QueryRowContext(
		ctx,
		query,
		category.Name,
	).Scan(&category.ID)
}

// ValidateCategory validates the fields of a category
func ValidateCategory(v *validator.Validator, category *Category) {
	v.Check(validator.NotBlank(category.Name), "category_name", "A category name is required")
	v.Check(validator.MaxLength(category.Name, 50), "category_name", "category name must be less than 50 characters")
	v.Check(validator.MinLength(category.Name, 3), "category_name", "category name must be at least 3 characters")
}

// GetAll retrieves all categories from the database
func (c *CategoryModel) GetAll() ([]*Category, error) {
	query := `SELECT id, name 
			  FROM categories 
			  ORDER BY name ASC`
	rows, err := c.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*Category{}
	for rows.Next() {
		category := &Category{}
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// Delete deletes a category from the database
func (c *CategoryModel) Delete(id int64) error {
	query := `DELETE FROM categories WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := c.DB.ExecContext(ctx, query, id)
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
