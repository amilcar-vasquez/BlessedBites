package main

import (
	"net/http"

	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/amilcar-vasquez/blessed-bites/internal/validator"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

// GET /signup handler to render the signup form
func (app *application) signupForm(w http.ResponseWriter, r *http.Request) {
	// Create a new template data instance
	data := NewTemplateData()
	data.Title = "Sign Up"
	data.HeaderText = "Sign Up"
	// Render the signup form template
	err := app.render(w, http.StatusOK, "signup.tmpl", data)
	if err != nil {
		app.logger.Error("Error rendering template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("Error parsing form data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Extract form fields
	email := r.PostForm.Get("email")
	fullname := r.PostForm.Get("fullname")
	phoneNo := r.PostForm.Get("phoneNo")
	password := r.PostForm.Get("password")

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		app.logger.Error("Error hashing password", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create an instance of user
	user := &data.User{
		Email:    email,
		FullName: fullname,
		PhoneNo:  phoneNo,
		Password: string(hashedPassword),
		Role:     "user", // default for now, dynamic later
	}

	// Validate the user data
	v := validator.NewValidator()
	data.ValidateUser(v, user)
	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Sign Up"
		data.HeaderText = "Sign Up"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"email":    email,
			"fullname": fullname,
			"phoneNo":  phoneNo,
			"password": password,
		}
		// Re-render the form with validation errors
		err := app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		if err != nil {
			app.logger.Error("Error rendering template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Insert the user into the database
	err = app.User.Insert(user)
	if err != nil {
		app.logger.Error("Error inserting user into database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//redirect to the thank you page
	http.Redirect(w, r, "/signup-thanks", http.StatusSeeOther)
}

// handler for signup thanks page
func (app *application) signupThanks(w http.ResponseWriter, r *http.Request) {
	// Create a new template data instance
	data := NewTemplateData()
	data.Title = "Thank You"
	data.HeaderText = "Thank You for Signing Up"
	// Render the thank you template
	err := app.render(w, http.StatusOK, "signup-thanks.tmpl", data)
	if err != nil {
		app.logger.Error("Error rendering template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// handler to delete a user
func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	userIDStr := r.PostForm.Get("user_id")
	if userIDStr == "" {
		app.logger.Error("user ID not provided")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Convert user ID to int64
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		app.logger.Error("invalid user ID format", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Delete the user from the database
	err = app.User.Delete(userID)
	if err != nil {
		app.logger.Error("failed to delete user", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Redirect to the user list page
	http.Redirect(w, r, "/user", http.StatusSeeOther)
}

// handler for rendering the login form
func (app *application) loginForm(w http.ResponseWriter, r *http.Request) {
	// Create a new template data instance
	data := NewTemplateData()
	data.Title = "Login"
	data.HeaderText = "Login"
	// Render the login form template
	err := app.render(w, http.StatusOK, "login.tmpl", data)
	if err != nil {
		app.logger.Error("Error rendering template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.User.GetByEmail(email)
	if err != nil {
		http.Error(w, "Email not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Store user ID in session
	session, _ := app.sessionStore.Get(r, "session")
	session.Values["userID"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/menu", http.StatusSeeOther)
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.sessionStore.Get(r, "session")
	delete(session.Values, "userID")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// handler to render the user page
func (app *application) userPageHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.User.GetAll()
	if err != nil {
		app.logger.Error("Error retrieving users", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Create a new template data instance
	data := NewTemplateData()
	data.Title = "User Page"
	data.HeaderText = "User Page"
	data.Users = users
	// Render the user page template
	err = app.render(w, http.StatusOK, "users.tmpl", data)
	if err != nil {
		app.logger.Error("Error rendering template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
