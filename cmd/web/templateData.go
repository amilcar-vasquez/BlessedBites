// file: cmd/web/templateData.go
package main

import (
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
)

type TemplateData struct {
	Title           string
	HeaderText      string
	MenuItems       []*data.MenuItem
	MenuItem        *data.MenuItem
	Categories      []*data.Category
	CategoryMap     map[int]string
	RandomMenuItems []*data.MenuItem
	Users           []*data.User
	User            *data.User
	IsAuthenticated bool
	CurrentUserID   int64
	CurrentUserRole string
	AlertMessage    string // To hold general messages like "Invalid credentials"
	AlertType       string // e.g., "alert-danger", "alert-success"
	CSRFField       template.HTML

	FormErrors map[string]string
	FormData   map[string]string
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

func (app *application) addDefaultData(td *TemplateData, r *http.Request) *TemplateData {
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

	// Handle flash messages (optional but common)
	if flash, ok := session.Values["flash"].(string); ok {
		td.AlertMessage = flash
		td.AlertType = "alert-success"
		delete(session.Values, "flash")
		_ = session.Save(r, nil)
	}

	return td
}
