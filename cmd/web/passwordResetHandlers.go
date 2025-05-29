// file: cmd/web/passwordResetHandlers.go
package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/amilcar-vasquez/blessed-bites/internal/validator"
)

// GET: Show form to enter email for reset
func (app *application) showPasswordResetRequestForm(w http.ResponseWriter, r *http.Request) {
	data := app.addDefaultData(NewTemplateData(), w, r)
	data.Title = "Forgot Password"
	app.render(w, http.StatusOK, "ResetPasswordRequest.tmpl", data)
}

// POST: Handle form submission (email input)
func (app *application) handlePasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	v := validator.NewValidator()

	data.ValidateEmail(v, email)
	for k, vErr := range v.Errors {
		app.logger.Error("Validation error", "field", k, "error", vErr)
	}

	token, err := app.User.InitiatePasswordReset(email)
	if err != nil {
		app.logger.Error("Reset request failed", "error", err)
		http.Error(w, "Email not found", http.StatusNotFound)
		return
	}

	// call the email service to send the reset link
	to := email
	subject := "Password Reset Request"
	body := fmt.Sprintf("To reset your password, please click the following link: \n\n%s/reset-password?token=%s\n\nIf you did not request this, please ignore this email.", *app.addr, token)
	err = app.mailer.Send(to, subject, body)
	if err != nil {
		app.logger.Error("Failed to send reset email", "error", err)
		http.Error(w, "Failed to send reset email", http.StatusInternalServerError)
		return
	}
	// Log the email sent
	app.logger.Info("Password reset email sent", "to", to)
	// Prepare data for success page
	data := app.addDefaultData(NewTemplateData(), w, r)
	data.Title = "Reset Email Sent"

	// Simulate success page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GET: Show reset form
func (app *application) showResetPasswordForm(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	data := app.addDefaultData(NewTemplateData(), w, r)
	data.Title = "Set New Password"
	data.Token = token
	app.render(w, http.StatusOK, "ResetPassword.tmpl", data)
}

// POST: Handle new password
func (app *application) handleResetPasswordSubmission(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	password := r.FormValue("password")
	confirm := r.FormValue("confirmPassword")

	v := validator.NewValidator()
	data.ValidatePassword(v, password)
	for k, vErr := range v.Errors {
		app.logger.Error("Validation error", "field", k, "error", vErr)
	}
	//check that passwords match
	if password != confirm {
		v.AddError("confirmPassword", "Passwords do not match")
	}
	if len(v.Errors) > 0 {
		data := app.addDefaultData(NewTemplateData(), w, r)
		data.Title = "Set New Password"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"token":           token,
			"password":        password,
			"confirmPassword": confirm,
		}
		err := app.render(w, http.StatusUnprocessableEntity, "ResetPassword.tmpl", data)
		if err != nil {
			app.logger.Error("Error rendering template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return

		}
	}

	err := app.User.FinalizePasswordReset(token, password)
	if err != nil {
		app.logger.Error("Finalize reset failed", "error", err)

		if strings.Contains(err.Error(), "invalid token") {
			http.Error(w, "Your reset link has expired.  Please request a new one.", http.StatusBadRequest)
		} else {
			http.Error(w, "Invalid reset token", http.StatusBadRequest)
		}
		return

	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
