// file: cmd/web/server.go
package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/amilcar-vasquez/blessed-bites/internal/csrf"
)

func (app *application) serve() error {
	// Use the SESSION_KEY for CSRF to ensure consistency with session store
	csrfKey := []byte(os.Getenv("SESSION_KEY"))
	if len(csrfKey) == 0 {
		csrfKey = []byte("session-key-32-bytes-long-different")
	}

	// Set secure to true if the app environment is NOT development
	secure := os.Getenv("APP_ENV") != "development"
	trusted := []string{"https://blessedbites.bz"}
	if os.Getenv("APP_ENV") == "development" {
		secure = false
		// allow local dev origins
		trusted = []string{"http://localhost:4000", "http://127.0.0.1:4000"}
	}

	// Log CSRF configuration for debugging
	if os.Getenv("APP_ENV") == "development" {
		app.logger.Info("CSRF Configuration",
			"secure", secure,
			"trusted_origins", trusted,
			"key_length", len(csrfKey))
	}

	// Store CSRF key in app for use in handlers
	app.csrfKey = csrfKey

	// Create custom CSRF middleware
	csrfConfig := csrf.Config{
		Key:    csrfKey,
		Secure: secure,
		MaxAge: 3600, // 1 hour
	}

	csrfMiddleware := csrf.Middleware(csrfConfig)

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
