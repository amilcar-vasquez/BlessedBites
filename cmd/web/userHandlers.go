package main

import (
	"net/http"

	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/amilcar-vasquez/blessed-bites/internal/validator"
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
	//extract form fields

	email := r.PostForm.Get("email")
	fullname := r.PostForm.Get("fullname")
	phoneNo := r.PostForm.Get("phoneNo")
	password := r.PostForm.Get("password")

	//Create an instance of user
	user := &data.User{
		Email:    email,
		FullName: fullname,
		PhoneNo:  phoneNo,
		Password: password,
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
		//re-render the form with validation errors
		err := app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		if err != nil {
			app.logger.Error("Error rendering template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	err = app.Users.Insert(user)
	if err != nil {
		app.logger.Error("Error inserting user into database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
