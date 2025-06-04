package main

import (
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/amilcar-vasquez/blessed-bites/internal/validator"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

// parseForm parses the form data from the request both for add and update
func parseUserForm(r *http.Request, isUpdate bool) (*data.User, map[string]string, map[string]string, error) {
	var formErrors = make(map[string]string)
	var formData = make(map[string]string)
	var role string

	//parse form
	err := r.ParseForm()
	if err != nil {
		formErrors["form"] = "Error parsing form data"
		return nil, formErrors, formData, err
	}

	if isUpdate {
		role = r.PostForm.Get("role")
		if role == "" {
			role = "user" // default role
		}
	} else {
		role = "user" // default role for new users
	}

	// Extract form fields
	idStr := r.PostForm.Get("user_id") // user ID for update
	fullname := r.PostForm.Get("fullname")
	email := r.PostForm.Get("email")
	phoneNo := r.PostForm.Get("phoneNo")
	password := r.PostForm.Get("password")
	confirmPassword := r.PostForm.Get("confirmPassword")

	//save raw form values
	formData["user_id"] = idStr
	formData["fullname"] = fullname
	formData["email"] = email
	formData["phoneNo"] = phoneNo
	formData["password"] = password
	formData["confirmPassword"] = confirmPassword

	// Convert user ID to int64
	var userID int64
	if idStr != "" {
		userID, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			formErrors["user_id"] = "Invalid user ID format"
			return nil, formErrors, formData, err
		}
	}

	// Create an instance of user
	user := &data.User{
		ID:       userID,
		FullName: fullname,
		Email:    email,
		PhoneNo:  phoneNo,
		Password: password,
		Role:     role,
	}

	return user, formErrors, formData, nil
}

// GET /signup handler to render the signup form
func (app *application) signupForm(w http.ResponseWriter, r *http.Request) {
	// Create a new template data instance
	data := app.addDefaultData(NewTemplateData(), w, r)
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
	// Parse form data using parseUserForm
	user, formErrors, formData, err := parseUserForm(r, false)
	if err != nil {
		app.logger.Error("Error parsing form data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Validate the user data
	v := validator.NewValidator()
	data.ValidateUser(v, user)
	for k, vErr := range v.Errors {
		formErrors[k] = vErr
	}

	//check that passwords match

	if formData["password"] != formData["confirmPassword"] {
		formErrors["password"] = "Passwords do not match"
	}

	// Check for validation errors
	if len(formErrors) > 0 {
		data := app.addDefaultData(NewTemplateData(), w, r)
		data.Title = "Sign Up"
		data.HeaderText = "Sign Up"
		data.FormErrors = formErrors
		data.FormData = formData
		// Re-render the form with validation errors
		err := app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		if err != nil {
			app.logger.Error("Error rendering template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		app.logger.Error("Error hashing password", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Insert the user into the database
	err = app.User.Insert(user)
	if err != nil {
		app.logger.Error("Error inserting user into database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect to the thank you page
	http.Redirect(w, r, "/signup-thanks", http.StatusSeeOther)
}

// handler for signup thanks page
func (app *application) signupThanks(w http.ResponseWriter, r *http.Request) {
	// Create a new template data instance
	data := app.addDefaultData(NewTemplateData(), w, r)
	data.Title = "Thank You"
	data.HeaderText = "Thank You for Signing Up"
	// Render the thank you template
	err := app.render(w, http.StatusOK, "signupThanks.tmpl", data)
	if err != nil {
		app.logger.Error("Error rendering template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// handler to render the update user form
func (app *application) updateUserForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	idStr := r.FormValue("user_id")
	// Check if ID is present in the form data
	if idStr == "" {
		app.logger.Error("User ID is missing in form data")
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.logger.Error("Invalid user ID", "value", idStr, "error", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := app.User.GetByID(int64(id))
	if err != nil {
		app.logger.Error("User not found", "error", err)
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}

	data := app.addDefaultData(NewTemplateData(), w, r)
	data.Title = "Update User"
	data.HeaderText = "Update User"
	data.User = user

	err = app.render(w, http.StatusOK, "signup.tmpl", data)
	if err != nil {
		app.logger.Error("Error rendering template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// handler to update a user
func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {
	user, formErrors, formData, err := parseUserForm(r, true)
	if err != nil {
		app.logger.Error("Error parsing form", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get ID separately
	idStr := formData["user_id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.logger.Error("Invalid user ID", "value", idStr, "error", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	user.ID = int64(id)

	v := validator.NewValidator()
	data.ValidateUser(v, user)
	for k, vErr := range v.Errors {
		formErrors[k] = vErr
	}

	if len(formErrors) > 0 {
		data := app.addDefaultData(NewTemplateData(), w, r)
		data.Title = "Update User"
		data.HeaderText = "Update User"
		data.FormErrors = formErrors
		data.FormData = formData
		data.User = user

		err = app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		if err != nil {
			app.logger.Error("Error rendering template with errors", "error", err)
		}
		return
	}

	err = app.User.Update(user)
	if err != nil {
		app.logger.Error("Error updating user", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
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

// handler to render the user page
func (app *application) userPageHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.User.GetAll()
	if err != nil {
		app.logger.Error("Error retrieving users", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Create a new template data instance
	data := app.addDefaultData(NewTemplateData(), w, r)
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
