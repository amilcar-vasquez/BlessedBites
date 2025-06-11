// cmd/web/orderHandlers.go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/amilcar-vasquez/blessed-bites/internal/utils"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"net/http"
	"os"
	"strings"
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

	// Parse request input (JSON or form)
	if isJSON {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON input", http.StatusBadRequest)
			return
		}
		items = input.Items
	} else {
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

	// Get session
	session, _ := app.sessionStore.Get(r, "session")

	// --- User detection logic ---
	var user *data.User

	// 1. Check for logged-in user in session
	user = app.contextGetUser(r)
	if user != nil {
		app.logger.Info("User detected from session", "userID", user.ID, "fullName", user.FullName)
	} else {
		app.logger.Info("No user detected in session")
	}

	// 2. Fallback to guest user via phone number
	if user == nil && input.PhoneNo != "" {
		existingUser, err := app.User.GetByPhone(input.PhoneNo)
		if err == nil && existingUser != nil {
			user = existingUser
		} else {
			user, err = app.User.CreateGuestUser(input.FullName, input.PhoneNo)
			if err != nil {
				app.logger.Error("Unable to create guest user", "error", err, "fullName", input.FullName, "phoneNo", input.PhoneNo)
				http.Error(w, "Unable to create guest user", http.StatusInternalServerError)
				return
			}
		}
	}

	// 3. Handle admin walk-in order override
	if user != nil && user.Role == "admin" {
		walkInFullName := r.FormValue("walkInFullName")
		walkInPhoneNo := r.FormValue("walkInPhoneNo")
		if walkInFullName == "" {
			app.logger.Error("Walk-in customer creation failed: missing full name", "adminID", user.ID)
			http.Error(w, "Full name required for walk-in", http.StatusBadRequest)
			return
		}
		if walkInPhoneNo == "" {
			//create randome 7 digits for phone number
			walkInPhoneNo = utils.RandomPhone()
		}
		walkInUser, err := app.User.CreateWalkInCustomer(walkInFullName, walkInPhoneNo)
		if err != nil {
			app.logger.Error("Failed to create walk-in customer", "error", err, "walkInFullName", walkInFullName, "adminID", user.ID)
			http.Error(w, "Failed to create walk-in customer", http.StatusInternalServerError)
			return
		}
		app.logger.Info("Walk-in customer created", "userID", walkInUser.ID, "fullName", walkInFullName)
		user = walkInUser
	}

	// Final guard: must have a user by this point
	if user == nil {
		http.Error(w, "User could not be determined", http.StatusInternalServerError)
		return
	}

	// --- Order validation ---
	if len(items) == 0 {
		session.AddFlash("Please add items before ordering.", "error")
		_ = session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// --- Transaction starts ---
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

	// Insert order and items
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

	// --- Async: WhatsApp Notification ---
	go func() {
		fullName := user.FullName
		if name, ok := session.Values["fullName"].(string); ok && name != "" {
			fullName = name
		}

		phoneNo := user.PhoneNo
		if sessionPhone, ok := session.Values["phoneNo"].(string); ok && sessionPhone != "" {
			phoneNo = sessionPhone
		}
		if len(phoneNo) > 0 && !strings.HasPrefix(phoneNo, "+501") {
			phoneNo = "+501" + phoneNo
		}

		sendWhatsAppMessages(phoneNo, messageItems, totalCost, fullName)
	}()

	// --- Async: Update Popular Items ---
	go func() {
		if err := app.MenuItem.UpdatePopularItems(); err != nil {
			app.logger.Error("Popular items update failed", "error", err)
		}
	}()

	// --- Response handling ---
	if isJSON {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Order placed successfully!",
			"note":    "Create an account later to track orders.",
		})
	} else {
		session.AddFlash(fmt.Sprintf("Order placed! Total: $%.2f", totalCost), "success")
		_ = session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
