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
	RandomMenuItems     []*data.MenuItem
	TopUserMenuItems    []*data.MenuItem
	Users               []*data.User
	User                *data.User
	Rating              []*data.Rating
	Recommendation      *data.Recommendation
	IsAuthenticated     bool
	CurrentUserID       int64
	CurrentUserRole     string
	CurrentUserFullName string
	AlertMessage        string // To hold general messages like "Invalid credentials"
	AlertType           string // e.g., "alert-danger", "alert-success"
	CSRFField           template.HTML
	FormErrors          map[string]string
	FormData            map[string]string
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
		td.CurrentUserFullName = fullName
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
