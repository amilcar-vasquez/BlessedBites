// cmd/web/orderHandlers.go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"net/http"
	"os"
)

type OrderItemInput struct {
	MenuItemID int     `json:"id"`
	ItemAmount int     `json:"qty"` // note: match JS key
	ItemName   string  `json:"name"`
	ItemPrice  float64 `json:"price"`
}

type WhatsAppOrderItem struct {
	Name     string
	Quantity int
	Subtotal float64
}

// create a function to call WhatsApp API and send the order details
func sendWhatsAppMessage(orderID int, items []WhatsAppOrderItem) error {
	message := fmt.Sprintf("🧾 New Order #%d\n", orderID)
	for _, item := range items {
		message += fmt.Sprintf("Item: %s\n", item.Name)
		message += fmt.Sprintf("Quantity: %d\n", item.Quantity)
		message += fmt.Sprintf("Subtotal: $%.2f\n", item.Subtotal)
		message += "------------------------\n"
	}
	message += fmt.Sprintf("\nTotal: $%.2f", orderTotal(items))

	// Twilio credentials (recommended: store in environment variables)
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	from := "whatsapp:+14155238886" // Twilio sandbox number
	to := os.Getenv("WHATSAPP_TO")  // e.g. "whatsapp:+5016082424"

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(message)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Failed to send WhatsApp message:", err.Error())
		return err
	}

	fmt.Println("✅ WhatsApp message sent!")
	return nil
}

func orderTotal(items []WhatsAppOrderItem) float64 {
	var total float64
	for _, item := range items {
		total += item.Subtotal
	}
	return total
}

func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	user := app.contextGetUser(r)
	if user == nil {
		http.Error(w, "User must be logged in to place an order.", http.StatusUnauthorized)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get order data from the form
	orderDataStr := r.FormValue("orderData")

if orderDataStr == "" || orderDataStr == "[]" {
	session, err := app.sessionStore.Get(r, "session")
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	session.AddFlash("Please add items to your order before proceeding.", "error")
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Session save error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

	// Check if the items slice is empty
	var items []OrderItemInput
if err := json.Unmarshal([]byte(orderDataStr), &items); err != nil {
	http.Error(w, "Invalid order data", http.StatusBadRequest)
	return
}

if len(items) == 0 {
	session, err := app.sessionStore.Get(r, "session")
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	session.AddFlash("Please add some yumminess to your order before proceeding.", "error")
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Session save error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	var messageItems []WhatsAppOrderItem

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

		messageItems = append(messageItems, WhatsAppOrderItem{
			Name:     item.ItemName,
			Quantity: item.ItemAmount,
			Subtotal: subtotal,
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

	// 🔔 Send WhatsApp notification
	go func(orderID int, messageItems []WhatsAppOrderItem) {
		if err := sendWhatsAppMessage(orderID, messageItems); err != nil {
			app.logger.Error("Failed to send WhatsApp notification", "error", err)
		}
	}(orderID, messageItems)

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
