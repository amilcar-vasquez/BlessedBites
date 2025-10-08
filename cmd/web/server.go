// file: cmd/web/server.go
package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/csrf"
)

func (app *application) serve() error {
	csrfKey := []byte("ZQnXOK/iAwl+wMHKrQxS1VEw+9KAZUq=")

	// Allow toggling secure cookies and trusted origins for development
	secure := true
	trusted := []string{"https://blessedbites.bz"}
	if os.Getenv("APP_ENV") == "development" {
		secure = false
		// allow local dev origins
		trusted = []string{"http://localhost:4000", "http://127.0.0.1:4000"}
	}

	csrfMiddleware := csrf.Protect(
		csrfKey,
		csrf.Secure(secure),
		csrf.SameSite(csrf.SameSiteDefaultMode),
		csrf.HttpOnly(true),
		csrf.Path("/"),
		csrf.TrustedOrigins(trusted),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.logger.Error("CSRF failure",
				"method", r.Method,
				"path", r.URL.Path,
				"form_token", r.FormValue("gorilla.csrf.Token"),
				"header_token", r.Header.Get("X-CSRF-Token"),
				"cookie", r.Header.Get("Cookie"),
				"origin", r.Header.Get("Origin"),
				"referer", r.Referer(),
			)
			http.Error(w, "Forbidden - CSRF token invalid", http.StatusForbidden)
		},
		),
		),
	)

	srv := &http.Server{
		Addr:         *app.addr,
		Handler:      csrfMiddleware(app.routes()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}
	app.logger.Info("starting server", "addr", srv.Addr, "handler", srv.Handler)
	return srv.ListenAndServe()
}
