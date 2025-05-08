// cmd/web/orderHandlers.go
package main

import (
	"encoding/json"
	"net/http"
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"fmt"
)

type OrderItemInput struct {
	MenuItemID int     `json:"id"`
	ItemAmount int     `json:"qty"`   // note: match JS key
	ItemName   string  `json:"name"`
	ItemPrice  float64 `json:"price"`
}


func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user == nil {
		http.Error(w, "User must be logged in to place an order.", http.StatusUnauthorized)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	orderDataStr := r.FormValue("orderData")
	var items []OrderItemInput
	if err := json.Unmarshal([]byte(orderDataStr), &items); err != nil {
		http.Error(w, "Invalid order data", http.StatusBadRequest)
		return
	}

	// Start transaction
	tx, err := app.Order.DB.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var totalCost float64
	var fullItems []data.OrderItem

	// Validate items and compute total
	for _, item := range items {
		menuItem, err := app.MenuItem.Get(int64(item.MenuItemID))
		if err != nil {
			http.Error(w, "Menu item not found", http.StatusNotFound)
			return
		}
		subtotal := menuItem.Price * float64(item.ItemAmount)
		totalCost += subtotal

		fullItems = append(fullItems, data.OrderItem{
			MenuItemID: item.MenuItemID,
			Quantity:   item.ItemAmount,
			ItemPrice:  menuItem.Price,
			Subtotal:   subtotal,
		})
	}

	// Insert into orders table
	orderID, err := app.Order.Insert(int(user.ID), totalCost)
	if err != nil {
		http.Error(w, "Could not save order", http.StatusInternalServerError)
		return
	}

	// Insert into order_items table
	for _, item := range fullItems {
		item.OrderID = orderID
		err := app.OrderItem.Insert(tx, item)
		if err != nil {
			http.Error(w, "Failed to save order items", http.StatusInternalServerError)
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, "Could not finalize order", http.StatusInternalServerError)
		return
	}

	// Send success flash with the total included
	session, err := app.sessionStore.Get(r, "session")
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	session.AddFlash(fmt.Sprintf("Order placed successfully! Total: $%.2f", totalCost), "success")
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Session save error", http.StatusInternalServerError)
		return
	}
	// Redirect to the order confirmation page or home page


	http.Redirect(w, r, "/", http.StatusSeeOther)
}
