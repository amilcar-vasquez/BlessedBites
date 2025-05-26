// data/orderitem.go
package data

import (
	"database/sql"
)

type OrderItem struct {
	ID         int
	OrderID    int
	MenuItemID int
	Quantity   int
	ItemPrice  float64
	Subtotal   float64
}

type OrderItemModel struct {
	DB *sql.DB
}

type SalesRecord struct {
	Date       string  `json:"Date"`
	ClientName string  `json:"ClientName"`
	Amount     float64 `json:"Amount"`
}

func (m OrderItemModel) Insert(tx *sql.Tx, item OrderItem) error {
	query := `
		INSERT INTO order_items (order_id, menu_item_id, quantity, item_price)
		VALUES ($1, $2, $3, $4)`
	_, err := tx.Exec(query, item.OrderID, item.MenuItemID, item.Quantity, item.ItemPrice)
	return err
}

// calculate daily sales for a specific date, returning a list of sales records with client name
func (m OrderItemModel) DailySales(date string) ([]SalesRecord, error) {
	query := `
		SELECT DATE(o.created_at) AS date, u.full_name AS client_name, SUM(oi.subtotal) AS amount
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		JOIN users u ON o.user_id = u.id
		WHERE DATE(o.created_at) = $1
		GROUP BY DATE(o.created_at), u.full_name
		ORDER BY u.full_name`
	rows, err := m.DB.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []SalesRecord
	for rows.Next() {
		var record SalesRecord
		if err := rows.Scan(&record.Date, &record.ClientName, &record.Amount); err != nil {
			return nil, err
		}
		sales = append(sales, record)
	}
	return sales, nil
}

// a slice of sales for the last 7 days
func (m OrderItemModel) Last7DaysSales() ([]SalesRecord, error) {
	query := `
		SELECT DATE(o.created_at) AS date, SUM(oi.subtotal)
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		WHERE o.created_at >= NOW() - INTERVAL '7 days'
		GROUP BY DATE(o.created_at)
		ORDER BY DATE(o.created_at)`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []SalesRecord
	for rows.Next() {
		var record SalesRecord
		if err := rows.Scan(&record.Date, &record.Amount); err != nil {
			return nil, err
		}
		sales = append(sales, record)
	}
	return sales, nil
}
