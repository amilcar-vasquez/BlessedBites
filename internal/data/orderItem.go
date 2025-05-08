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

func (m OrderItemModel) Insert(tx *sql.Tx, item OrderItem) error {
	query := `
		INSERT INTO order_items (order_id, menu_item_id, quantity, item_price)
		VALUES ($1, $2, $3, $4)`
	_, err := tx.Exec(query, item.OrderID, item.MenuItemID, item.Quantity, item.ItemPrice)
	return err
}
