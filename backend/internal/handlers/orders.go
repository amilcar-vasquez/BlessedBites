package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/amilcar-vasquez/blessed-bites/backend/internal/middleware"
	"github.com/amilcar-vasquez/blessed-bites/backend/internal/realtime"
	"github.com/amilcar-vasquez/blessed-bites/backend/pkg/responses"
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
)

type OrdersHandler struct {
	Orders    *data.OrderModel
	OrderItem *data.OrderItemModel
	MenuItem  *data.MenuItemModel
	User      *data.UserModel
	Broker    *realtime.Broker
}

type OrderItemInput struct {
	MenuItemID int `json:"id"`
	ItemAmount int `json:"qty"`
}

type OrderRequest struct {
	Items    []OrderItemInput `json:"items"`
	FullName string           `json:"full_name"`
	PhoneNo  string           `json:"phone_no"`
}

func NewOrdersHandler(db *sql.DB, broker *realtime.Broker) *OrdersHandler {
	return &OrdersHandler{
		Orders:    &data.OrderModel{DB: db},
		OrderItem: &data.OrderItemModel{DB: db},
		MenuItem:  &data.MenuItemModel{DB: db},
		User:      &data.UserModel{DB: db},
		Broker:    broker,
	}
}

func (h *OrdersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequest(w, "invalid JSON payload")
		return
	}
	if len(input.Items) == 0 {
		responses.BadRequest(w, "order items are required")
		return
	}

	user, err := h.resolveUser(r, input.FullName, input.PhoneNo)
	if err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	tx, err := h.Orders.DB.Begin()
	if err != nil {
		responses.InternalServerError(w, "failed to start transaction")
		return
	}
	defer tx.Rollback()

	var (
		orderItems []data.OrderItem
		totalCost  float64
	)
	for _, item := range input.Items {
		menuItem, err := h.MenuItem.Get(int64(item.MenuItemID))
		if err != nil {
			responses.NotFound(w, "menu item not found")
			return
		}

		subtotal := float64(item.ItemAmount) * menuItem.Price
		totalCost += subtotal
		orderItems = append(orderItems, data.OrderItem{
			MenuItemID: item.MenuItemID,
			Quantity:   item.ItemAmount,
			ItemPrice:  menuItem.Price,
			Subtotal:   subtotal,
		})
	}

	orderID, err := h.Orders.Insert(int(user.ID), totalCost)
	if err != nil {
		responses.InternalServerError(w, "failed to create order")
		return
	}

	for _, item := range orderItems {
		item.OrderID = orderID
		if err := h.OrderItem.Insert(tx, item); err != nil {
			responses.InternalServerError(w, "failed to add order item")
			return
		}
		_ = h.MenuItem.IncrementOrderCount(int64(item.MenuItemID))
	}

	if err := tx.Commit(); err != nil {
		responses.InternalServerError(w, "failed to commit order")
		return
	}

	go func() {
		_ = h.MenuItem.UpdatePopularItems()
		h.Broker.Publish(realtime.Event{
			Type: "order.created",
			Data: map[string]any{
				"order_id":   orderID,
				"user_id":    user.ID,
				"total_cost": totalCost,
			},
		})
	}()

	responses.JSON(w, http.StatusCreated, map[string]any{
		"order_id": orderID,
		"message":  "order placed successfully",
	})
}

func (h *OrdersHandler) resolveUser(r *http.Request, fullName, phone string) (*data.User, error) {
	if claims, ok := middleware.ClaimsFromContext(r.Context()); ok {
		user, err := h.User.GetByID(claims.UserID)
		if err == nil && user != nil {
			return user, nil
		}
	}

	fullName = strings.TrimSpace(fullName)
	phone = strings.TrimSpace(phone)
	if fullName == "" || phone == "" {
		return nil, errBadGuestPayload
	}

	user, err := h.User.GetByPhone(phone)
	if err == nil && user != nil {
		return user, nil
	}

	return h.User.CreateGuestUser(fullName, phone)
}

var errBadGuestPayload = &apiError{message: "guest checkout requires full_name and phone_no"}

type apiError struct {
	message string
}

func (e *apiError) Error() string {
	return e.message
}
