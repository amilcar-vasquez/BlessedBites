// file: cmd/web/debugHandlers.go
package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
)

// debugCSRFHandler returns the server-side CSRF token for the current request
// and echoes back the CSRF cookie value for comparison. This handler is only
// intended for local development and is gated by APP_ENV=development.
func (app *application) debugCSRFHandler(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("APP_ENV") != "development" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	token := csrf.Token(r)
	cookie, _ := r.Cookie("_gorilla_csrf")
	cookieVal := ""
	if cookie != nil {
		cookieVal = cookie.Value
	}

	out := map[string]string{
		"server_token": token,
		"cookie_value": cookieVal,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
