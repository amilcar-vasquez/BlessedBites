package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

// getSessionSafe tries to retrieve a session and, if the cookie value is invalid
// (securecookie decoding error), returns a new empty session instead of failing.
func (app *application) getSessionSafe(w http.ResponseWriter, r *http.Request, name string) (*sessions.Session, error) {
	session, err := app.sessionStore.Get(r, name)
	if err == nil {
		return session, nil
	}

	// If it's the securecookie invalid error, return a new session without failing
	if strings.Contains(err.Error(), "securecookie: the value is not valid") {
		// log and clear the cookie so the browser replaces it
		app.logger.Info("Invalid session cookie detected; issuing new session and clearing cookie", "name", name)

		// create new session
		newSess, _ := app.sessionStore.New(r, name)

		// Clear cookie on response to remove invalid value (if we have a writer)
		// Build cookie with MaxAge= -1 to instruct browser to delete it
		if w != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     name,
				Value:    "",
				Path:     newSess.Options.Path,
				MaxAge:   -1,
				HttpOnly: newSess.Options.HttpOnly,
				Secure:   newSess.Options.Secure,
				SameSite: newSess.Options.SameSite,
				Domain:   newSess.Options.Domain,
			})
		}

		return newSess, nil
	}

	return session, err
}
