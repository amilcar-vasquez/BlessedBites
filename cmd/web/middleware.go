// file: cmd/web/middleware.go
package main

import (
	"fmt"
	"net/http"
)

func (app *application) loggingMiddleware(next http.Handler) http.Handler {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info(
			"received request",
			"ip", ip,
			"protocol", proto,
			"method", method,
			"uri", uri,
		)
		next.ServeHTTP(w, r)
		app.logger.Info("Request processed")
	})
	return fn

}

// development middleware to allow css and js to be reloaded
func noCacheMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, must-revalidate")
		w.Header().Set("Expires", "0")
		w.Header().Set("Pragma", "no-cache")
		h.ServeHTTP(w, r)
	})
}

// requireAuth is a middleware that checks if the user is authenticated.
func (app *application) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.getSessionSafe(w, r, "session")

		auth, ok := session.Values["authenticated"].(bool)

		// Debug logging
		app.logger.Info("requireAuth middleware",
			"path", r.URL.Path,
			"app.env", app.env,
			"auth", auth,
			"ok", ok,
			"session_values", session.Values)

		// ✅ Bypass auth in local development
		if app.env == "development" {
			app.logger.Info("Bypassing authentication for localhost (development mode)")
			session.Values["authenticated"] = true
			session.Values["authenticatedUserID"] = int64(1)
			session.Values["userRole"] = "admin"
			session.Values["fullName"] = "Local Dev"
			err := session.Save(r, w)
			if err != nil {
				app.logger.Error("Failed to save session in development bypass", "error", err)
			}
			next.ServeHTTP(w, r)
			return
		}

		if !ok || !auth {
			flashSession, _ := app.getSessionSafe(w, r, "flash")
			flashSession.Values["alertMessage"] = "This is an admin-only area. Please log in first."
			flashSession.Values["alertType"] = "alert-warning"
			flashSession.Save(r, w)

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (app *application) IsAuthenticated(r *http.Request) bool {
	session, err := app.getSessionSafe(nil, r, "session")
	if err != nil {
		app.logger.Error("Failed to get session in IsAuthenticated", "error", err)
		return false
	}
	if session == nil {
		return false
	}
	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}

func (app *application) CurrentUserID(r *http.Request) (int64, error) {
	session, err := app.getSessionSafe(nil, r, "session")
	if err != nil {
		return 0, err
	}
	if session == nil {
		return 0, fmt.Errorf("user ID not found in session")
	}
	// stored user id may be int or int64 depending on where it was set
	if id64, ok := session.Values["authenticatedUserID"].(int64); ok {
		return id64, nil
	}
	if idInt, ok := session.Values["authenticatedUserID"].(int); ok {
		return int64(idInt), nil
	}
	return 0, fmt.Errorf("user ID not found in session")
}

func (app *application) CurrentUserFullName(r *http.Request) (string, error) {
	session, err := app.getSessionSafe(nil, r, "session")
	if err != nil {
		return "", err
	}
	if session == nil {
		return "", fmt.Errorf("user full name not found in session")
	}
	fullName, ok := session.Values["fullName"].(string)
	if !ok {
		return "", fmt.Errorf("user full name not found in session")
	}
	return fullName, nil
}
