// data/categories.go
package data

import (
	"database/sql"
)

type Category struct {
	ID   int
	Name string
}

type CategoryModel struct {
	DB *sql.DB
}

func (c CategoryModel) GetAll() ([]Category, error) {
	query := `SELECT id, name FROM categories`
	rows, err := c.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func NewCategoryModel(db *sql.DB) CategoryModel {
	return CategoryModel{DB: db}
}
