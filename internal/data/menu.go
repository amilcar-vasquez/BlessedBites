// data/menu.go
package data

import (
	"database/sql"
)

type MenuItem struct {
	ID          int
	Name        string
	Siding      string
	Price       float64
	Description string
	CategoryID  int
	OrderCount  int
}

type MenuItemModel struct {
	DB *sql.DB
}

func (m MenuItemModel) GetAll() ([]MenuItem, error) {
	query := `SELECT id, name, siding, price, description, category_id, order_count FROM menu_items`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []MenuItem
	for rows.Next() {
		var item MenuItem
		err := rows.Scan(&item.ID, &item.Name, &item.Siding, &item.Price, &item.Description, &item.CategoryID, &item.OrderCount)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func NewMenuItemModel(db *sql.DB) MenuItemModel {
	return MenuItemModel{DB: db}
}
