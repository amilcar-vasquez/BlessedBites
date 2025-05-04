// File: cmd/web/loginHandlers.go
package main

import (
	"fmt"
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/amilcar-vasquez/blessed-bites/internal/validator"
	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"strings"
)

// handler for rendering the login form
func (app *application) loginForm(w http.ResponseWriter, r *http.Request) {
	//check if user is already logged in
	session, err := app.sessionStore.Get(r, "session")
	if err != nil {
		app.logger.Error("Failed to get session", "error", err)
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}
	isAuthenticated, ok := session.Values["authenticated"].(bool)
	if ok && isAuthenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//if not logged in, render the login form
	data := app.addDefaultData(NewTemplateData(), r)

	data.CSRFField = csrf.TemplateField(r)
	session, _ = app.sessionStore.Get(r, "signup-data")
	if session.Values["email"] != nil {
		alertMessage := "Sign up was successful with (" +
			"Email: " + session.Values["email"].(string) + " , " +
			"Password: " + session.Values["password"].(string) +
			")"

		// assign the string to the alert message
		data.AlertMessage = alertMessage
		data.AlertType = "success"

		session.Options.MaxAge = -1 // delete the session
		session.Save(r, w)
	}
	fmt.Println("CSRF token:", csrf.Token(r))

	err = app.render(w, http.StatusOK, "login.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render login page", "template", "signin.tmpl", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// handler for processing the login form
func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.Form.Get("email"))
	password := r.Form.Get("password")

	dummyHash := []byte("$2a$10$CwTycUXWue0Thq9StjUM0uJ7kYCcvl5yE9Ew2yHMTKJ2HZJY5t1L6") // bcrypt hash for "password"

	userData := &data.User{
		Email:    email,
		Password: password,
	}

	// Validate the email and password
	v := validator.NewValidator()
	data.ValidateLogin(v, userData)

	if !v.ValidData() {
		data := app.addDefaultData(NewTemplateData(), r)
		data.CSRFField = template.HTML(csrf.TemplateField(r))

		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"email": email,
		}
		data.AlertMessage = "Please correct the errors below."
		data.AlertType = "alert-warning"

		err := app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render user Form with validation errors", "template", "signin.tmpl", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		return // Stop processing
	}

	user, err := app.User.GetByEmail(email)

	if err != nil {
		data := app.addDefaultData(NewTemplateData(), r)
		data.CSRFField = template.HTML(csrf.TemplateField(r))
		bcrypt.CompareHashAndPassword(dummyHash, []byte(password)) // mitigate timing attack
		data.AlertMessage = "Invalid email or password."
		data.AlertType = "alert-danger"
		data.FormData = map[string]string{
			"email": email,
		}
		app.render(w, http.StatusUnauthorized, "login.tmpl", data)
		return
	}

	user = &data.User{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		PhoneNo:   user.PhoneNo,
		Password:  user.Password,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			data := app.addDefaultData(NewTemplateData(), r)
			data.CSRFField = template.HTML(csrf.TemplateField(r))
			data.AlertMessage = "Invalid email or password"
			data.AlertType = "danger"
			data.FormData = map[string]string{
				"email": email,
			}
			app.render(w, http.StatusUnauthorized, "login.tmpl", data)
		} else {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		}
		return
	}
	// Store user ID and role in session
	session, err := app.sessionStore.Get(r, "session")
	if err != nil {
		app.logger.Error("Failed to get session", "error", err)
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = true
	session.Values["authenticatedUserID"] = user.ID // <- store the user ID
	session.Values["userRole"] = user.Role          // <- store the role
	session.Options.MaxAge = 3600                   // Set session expiration to 1 hour

	//also set these values in the template data
	data := app.addDefaultData(NewTemplateData(), r)
	data.CSRFField = csrf.TemplateField(r)
	data.IsAuthenticated = true
	data.CurrentUserID = user.ID
	data.CurrentUserRole = user.Role
	data.AlertMessage = "Login successful"
	data.AlertType = "alert-success"

	// Save the session
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := app.sessionStore.Get(r, "session")
	if err != nil {
		app.logger.Error("Failed to get session", "error", err)
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}

	// Log the logout
	if id, ok := session.Values["authenticatedUserID"]; ok {
		app.logger.Info("User logged out", "userID", id)
	}

	// Invalidate session
	session.Options.MaxAge = -1

	// Optional: set a flash message before clearing
	// If you're not using flash messages, you can skip this
	flashSession, _ := app.sessionStore.New(r, "flash")
	flashSession.Values["alertMessage"] = "You have been logged out successfully."
	flashSession.Values["alertType"] = "alert-success"

	err = session.Save(r, w)
	if err != nil {
		app.logger.Error("failed to save invalidated session", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = flashSession.Save(r, w)
	if err != nil {
		app.logger.Error("failed to save flash message", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
