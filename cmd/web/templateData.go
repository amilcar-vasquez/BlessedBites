// file: cmd/web/templateData.go
package main

import (
	"html/template"
	"net/http"
	"os"

	"github.com/amilcar-vasquez/blessed-bites/internal/csrf"
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
)

type TemplateData struct {
	Title               string
	HeaderText          string
	MenuItems           []*data.MenuItem
	MenuItem            *data.MenuItem
	Categories          []*data.Category
	CategoryMap         map[int]string
	TopUserMenuItems    []*data.MenuItem
	Users               []*data.User
	User                *data.User
	Rating              []*data.Rating
	Last7DaysSales      []data.SalesRecord // Sales data for the last 7 days
	Recommendation      *data.Recommendation
	IsAuthenticated     bool
	CurrentUserID       int64
	CurrentUserRole     string
	CurrentUserFullName string
	CurrentUserPhone    string
	AlertMessage        string // To hold general messages like "Invalid credentials"
	AlertType           string // e.g., "alert-danger", "alert-success"
	CSRFField           template.HTML
	FormErrors          map[string]string
	FormData            map[string]string
	TotalOrders         int
	DailySales          []data.SalesRecord // Total sales for the day
	Top5MenuItems       []*data.MenuItem   // Top 5 popular menu items
	OrderItems          []*data.OrderItem  // Added from orderItem.go
	ChartLabels         []string
	ChartData           []float64
	Token               string // For password reset
	CurrentPage         int    // For pagination
	TotalPages          int    // For pagination
}

// factory function to initialize a new templateData struct
func NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:      "Welcome to Blessed Bites",
		HeaderText: "Welcome to Blessed Bites",
		FormErrors: map[string]string{},
		FormData:   map[string]string{},
	}
}

func (app *application) addDefaultData(td *TemplateData, w http.ResponseWriter, r *http.Request) *TemplateData {
	if td == nil {
		td = &TemplateData{}
	}

	// Get session for reading auth status
	// obtain session safely (don't fail on invalid cookie)
	session, _ := app.getSessionSafe(nil, r, "session")

	// Generate CSRF token using our custom implementation
	secure := os.Getenv("APP_ENV") != "development"
	td.CSRFField = csrf.TemplateField(r, app.csrfKey, w, secure)

	// Check authentication status WITHOUT modifying session
	if userID, exists := session.Values["authenticatedUserID"]; exists {
		td.IsAuthenticated = true

		if id, ok := userID.(int); ok {
			td.CurrentUserID = int64(id)
		}
	}

	// Debug: Log all session values
	if os.Getenv("APP_ENV") == "development" {
		app.logger.Info("Session values debug",
			"authenticatedUserID", session.Values["authenticatedUserID"],
			"userRole", session.Values["userRole"],
			"fullName", session.Values["fullName"],
			"phoneNo", session.Values["phoneNo"],
			"authenticated", session.Values["authenticated"])
	}

	// Get user details from session (avoid database calls)
	if role, exists := session.Values["userRole"]; exists {
		if roleStr, ok := role.(string); ok {
			td.CurrentUserRole = roleStr
		}
	}

	if fullName, exists := session.Values["fullName"]; exists {
		if nameStr, ok := fullName.(string); ok {
			td.CurrentUserFullName = nameStr
		}
	}

	if phone, exists := session.Values["phoneNo"]; exists {
		if phoneStr, ok := phone.(string); ok {
			td.CurrentUserPhone = phoneStr
		}
	}

	// Handle flash messages - READ ONLY, don't modify session
	if flash, exists := session.Values["flash"]; exists && flash != "" {
		if flashStr, ok := flash.(string); ok && flashStr != "" {
			td.AlertMessage = flashStr
			td.AlertType = "alert-info" // default
			// NOTE: Flash removal should be handled by individual handlers, not here
			// Modifying sessions in addDefaultData breaks CSRF token validation
		}
	}

	return td
}
