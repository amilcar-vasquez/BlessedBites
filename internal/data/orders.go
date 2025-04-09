// data/orders.go
package data

import (
	"database/sql"
)

type Order struct {
	ID         int
	MenuItemID int
	ItemAmount int
	MenuPrice  float64
	TotalCost  float64
	UserID     int
	CreatedAt  string // adjust to time.Time if using timestamps
}

type OrderModel struct {
	DB *sql.DB
}

func (o OrderModel) Insert(order Order) error {
	query := `INSERT INTO orders (menu_item_id, item_amount, menu_price, total_cost, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())`
	_, err := o.DB.Exec(query, order.MenuItemID, order.ItemAmount, order.MenuPrice, order.TotalCost, order.UserID)
	return err
}

func (o OrderModel) GetByUser(userID int) ([]Order, error) {
	query := `SELECT id, menu_item_id, item_amount, menu_price, total_cost, user_id, created_at FROM orders WHERE user_id = $1`
	rows, err := o.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.MenuItemID, &order.ItemAmount, &order.MenuPrice, &order.TotalCost, &order.UserID, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func NewOrderModel(db *sql.DB) OrderModel {
	return OrderModel{DB: db}
}
