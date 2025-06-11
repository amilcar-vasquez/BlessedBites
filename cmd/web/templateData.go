// file: cmd/web/templateData.go
package main

import (
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
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
	td.CSRFField = csrf.TemplateField(r)
	session, _ := app.sessionStore.Get(r, "session")

	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		td.IsAuthenticated = true
	}

	if role, ok := session.Values["userRole"].(string); ok {
		td.CurrentUserRole = role
	}

	if userID, ok := session.Values["authenticatedUserID"].(int64); ok {
		td.CurrentUserID = userID
	}

	if fullName, ok := session.Values["fullName"].(string); ok {
		// Set CurrentUserFullName to just the first word (first name)
		if len(fullName) > 0 {
			for i, r := range fullName {
				if r == ' ' {
					td.CurrentUserFullName = fullName[:i]
					break
				}
			}
			if td.CurrentUserFullName == "" {
				td.CurrentUserFullName = fullName
			}
		}
	}

	//add logic for supporting flash messages

	if flashes := session.Flashes("success"); len(flashes) > 0 {
		td.AlertMessage = flashes[0].(string)
		td.AlertType = "success" // NEW: Custom class instead of Materialize colors
	}
	if flashes := session.Flashes("error"); len(flashes) > 0 {
		td.AlertMessage = flashes[0].(string)
		td.AlertType = "error" // You can customize more types this way
	}
	if err := session.Save(r, w); err != nil {
		app.logger.Error("Error saving session", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil
	}
	return td

}
