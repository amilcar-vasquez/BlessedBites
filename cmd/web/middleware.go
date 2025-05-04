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

		app.logger.Info("received request", "ip", ip, "protocol", proto, "method", method, "uri", uri)
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
		session, _ := app.sessionStore.Get(r, "session")
		userID := session.Values["userID"]
		if userID == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// IsAuthenticated checks if the user is authenticated by checking the session.
func (app *application) IsAuthenticated(r *http.Request) bool {
	session, _ := app.sessionStore.Get(r, "session-name")
	_, ok := session.Values["userID"]
	return ok
}

// CurrentUserID retrieves the current user's ID from the session.
func (app *application) CurrentUserID(r *http.Request) (int64, error) {
	session, _ := app.sessionStore.Get(r, "session-name")
	id, ok := session.Values["userID"].(int64)
	if !ok {
		return 0, fmt.Errorf("user ID not found in session")
	}
	return id, nil
}
