// file: internal/data/order.go
package data

import (
	"database/sql"
	"time"
)

type Order struct {
	ID            int
	UserID        int
	TotalCost     float64
	CreatedAt     time.Time
	Status        string
	PaymentMethod sql.NullString
}

type OrderModel struct {
	DB *sql.DB
}

func (m OrderModel) Insert(userID int, total float64) (int, error) {
	var orderID int
	query := `
		INSERT INTO orders (user_id, total_cost)
		VALUES ($1, $2)
		RETURNING id`
	err := m.DB.QueryRow(query, userID, total).Scan(&orderID)
	return orderID, err
}

func (m OrderModel) GetByUser(userID int) ([]Order, error) {
	query := `SELECT id, user_id, total_cost, created_at, status, payment_method FROM orders WHERE user_id = $1`
	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.UserID, &order.TotalCost, &order.CreatedAt, &order.Status, &order.PaymentMethod)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
