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

type OrderRequest struct {
	Items    []OrderItemInput `json:"items"`
	FullName string           `json:"full_name"`
	PhoneNo  string           `json:"phone_no"`
	// Maybe optional email field for account conversion later
}

type WhatsAppOrderItem struct {
	Name     string
	Quantity int
	Subtotal float64
}

// create a function to call WhatsApp API and send the order details
func sendWhatsAppMessages(customerPhone string, items []WhatsAppOrderItem, total float64, customerName string) {
	client := twilio.NewRestClient()

	//build itemized message strings
	var itemLines string
	for _, item := range items {
		line := fmt.Sprintf("%s x%d - $%.2f", item.Name, item.Quantity, item.Subtotal)
		itemLines += line + "\n"
	}

	// ===== 1. Send Freeform Message to Kitchen =====
	kitchenMessage := fmt.Sprintf(
		"New order from %s (%s):\n%sTotal: $%.2f",
		customerName,
		customerPhone,
		itemLines,
		total,
	)

	kitchenParams := &openapi.CreateMessageParams{}
	kitchenParams.SetTo(os.Getenv("WHATSAPP_TO")) // Store your kitchen WhatsApp number as env var
	kitchenParams.SetFrom(os.Getenv("TWILIO_FROM_NUMBER"))
	kitchenParams.SetBody(kitchenMessage)

	_, err := client.Api.CreateMessage(kitchenParams)
	if err != nil {
		fmt.Printf("Failed to send kitchen message: %v", err.Error())
	} else {
		fmt.Println("Kitchen message sent successfully")
	}

	/* ===== 2. Send Template Message to Customer =====
	customerParams := &openapi.CreateMessageParams{}
	customerParams.SetTo(customerPhone)
	customerParams.SetFrom(from)
	customerParams.SetMessagingServiceSid(os.Getenv("TWILIO_MSG_SID"))
	customerParams.SetContentSid(os.Getenv("TWILIO_ORDER_TEMPLATE_SID")) // Store your template SID as env var

	// Match the order of variables in your template
	contentVars := fmt.Sprintf(`{
		"1": "%s",
		"2": "%s",
		"3": "%s",
		"4": "%s",
		"5": "%s"
	}`, quantityStr, subtotalStr, item, customerName, totalStr)

	customerParams.SetContentVariables(contentVars)

	_, err = client.Api.CreateMessage(customerParams)
	if err != nil {
		fmt.Printf("Failed to send customer message: %v", err.Error())
	} else {
		fmt.Println("Customer message sent successfully")
	}*/
}

func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var input OrderRequest
	var items []OrderItemInput

	contentType := r.Header.Get("Content-Type")
	isJSON := contentType == "application/json"

	// Parse JSON input (mobile/webapp)
	if isJSON {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON input", http.StatusBadRequest)
			return
		}
		items = input.Items
	} else {
		// Handle form submission (web form / admin POS)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		input.FullName = r.FormValue("guestUserName")
		input.PhoneNo = r.FormValue("guestUserPhone")

		orderDataStr := r.FormValue("orderData")
		if orderDataStr == "" || orderDataStr == "[]" {
			session, _ := app.sessionStore.Get(r, "session")
			session.AddFlash("Please add items to your order before proceeding.", "error")
			_ = session.Save(r, w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err := json.Unmarshal([]byte(orderDataStr), &items); err != nil {
			http.Error(w, "Invalid order JSON data", http.StatusBadRequest)
			return
		}
	}

	// Check for walk-in admin order
	var user *data.User
	existingUser, err := app.User.GetByPhone(input.PhoneNo)
	if err == nil && existingUser != nil {
		user = existingUser
	} else {
		user, err = app.User.CreateGuestUser(input.FullName, input.PhoneNo)
		if err != nil {
			http.Error(w, "Unable to create guest user", http.StatusInternalServerError)
			return
		}
	}

	// Admin placing walk-in order?
	if user.Role == "admin" {
		fullName := r.FormValue("walkInFullName")
		if fullName == "" {
			http.Error(w, "Full name required for walk-in", http.StatusBadRequest)
			return
		}
		walkInUser, err := app.User.CreateWalkInCustomer(fullName)
		if err != nil {
			http.Error(w, "Failed to create walk-in customer", http.StatusInternalServerError)
			return
		}
		app.logger.Info("Walk-in customer created", "userID", walkInUser.ID, "fullName", fullName)
		user = walkInUser
	}

	// Validate order items
	if len(items) == 0 {
		session, _ := app.sessionStore.Get(r, "session")
		session.AddFlash("Please add items before ordering.", "error")
		_ = session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Begin transaction
	tx, err := app.Order.DB.Begin()
	if err != nil {
		http.Error(w, "Transaction error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var fullItems []data.OrderItem
	var messageItems []WhatsAppOrderItem
	var totalCost float64

	for _, item := range items {
		menuItem, err := app.MenuItem.Get(int64(item.MenuItemID))
		if err != nil {
			http.Error(w, "Menu item not found", http.StatusNotFound)
			return
		}
		sub := float64(item.ItemAmount) * menuItem.Price
		totalCost += sub

		fullItems = append(fullItems, data.OrderItem{
			MenuItemID: item.MenuItemID,
			Quantity:   item.ItemAmount,
			ItemPrice:  menuItem.Price,
			Subtotal:   sub,
		})

		messageItems = append(messageItems, WhatsAppOrderItem{
			Name:     item.ItemName,
			Quantity: item.ItemAmount,
			Subtotal: sub,
		})
	}

	// Insert order
	orderID, err := app.Order.Insert(int(user.ID), totalCost)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	for _, item := range fullItems {
		item.OrderID = orderID
		if err := app.OrderItem.Insert(tx, item); err != nil {
			http.Error(w, "Could not insert order items", http.StatusInternalServerError)
			return
		}
		_ = app.MenuItem.IncrementOrderCount(int64(item.MenuItemID))
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Commit failed", http.StatusInternalServerError)
		return
	}

	// Async WhatsApp Notification
	go func() {
		// Try to get user full name from session if available
		fullName := user.FullName
		session, _ := app.sessionStore.Get(r, "session")
		if name, ok := session.Values["fullName"].(string); ok && name != "" {
			fullName = name
		}

		// Get phone number from session if available, else use user.PhoneNo
		phoneNo := user.PhoneNo
		if sessionPhone, ok := session.Values["phoneNo"].(string); ok && sessionPhone != "" {
			phoneNo = sessionPhone
		}

		// Add country code +501 if not already present
		if len(phoneNo) > 0 && phoneNo[:4] != "+501" {
			phoneNo = "+501" + phoneNo
		}

		sendWhatsAppMessages(
			phoneNo,
			messageItems,
			totalCost,
			fullName,
		)
	}()

	// Async popular update
	go func() {
		if err := app.MenuItem.UpdatePopularItems(); err != nil {
			app.logger.Error("Popular items update failed", "error", err)
		}
	}()

	if isJSON {
		// Mobile/web JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Order placed successfully!",
			"note":    "Create an account later to track orders.",
		})
	} else {
		// Web form success redirect
		session, _ := app.sessionStore.Get(r, "session")
		session.AddFlash(fmt.Sprintf("Order placed! Total: $%.2f", totalCost), "success")
		_ = session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
