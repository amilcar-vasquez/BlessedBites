// filepath: cmd/web/main.go
package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/amilcar-vasquez/blessed-bites/internal/mailer"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

// Initialize session store with proper configuration
func initSessionStore() *sessions.CookieStore {
	// Support SESSION_KEYS (comma-separated) or fallback to SESSION_KEY.
	// Values can be raw or base64-encoded. Multiple keys allow rotation.
	keysEnv := os.Getenv("SESSION_KEYS")
	var keys [][]byte

	if keysEnv != "" {
		parts := strings.Split(keysEnv, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if decoded, err := base64.StdEncoding.DecodeString(p); err == nil && len(decoded) > 0 {
				keys = append(keys, decoded)
			} else if p != "" {
				keys = append(keys, []byte(p))
			}
		}
	} else {
		k := os.Getenv("SESSION_KEY")
		if k == "" {
			k = "session-key-32-bytes-long-different"
		}
		if decoded, err := base64.StdEncoding.DecodeString(k); err == nil && len(decoded) > 0 {
			keys = append(keys, decoded)
		} else {
			keys = append(keys, []byte(k))
		}
	}

	store := sessions.NewCookieStore(keys...)

	// Configure session options with DIFFERENT cookie name to avoid CSRF conflicts
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   os.Getenv("APP_ENV") != "development", // Only secure in production
		SameSite: http.SameSiteDefaultMode,
		Domain:   "", // Explicit empty domain for localhost
	}

	return store
}

var sessionStore = initSessionStore()

type application struct {
	addr           *string
	DB             *sql.DB // <- Add this line
	MenuItem       *data.MenuItemModel
	Order          *data.OrderModel
	OrderItem      *data.OrderItemModel
	Category       *data.CategoryModel
	User           *data.UserModel
	sessionStore   *sessions.CookieStore
	Analytics      *data.AnalyticsModel
	Rating         *data.RatingModel
	Recommendation *data.RecommendationModel
	logger         *slog.Logger
	templateCache  map[string]*template.Template
	mailer         *mailer.Mailer // Add mailer field
	env            string
	csrfKey        []byte // CSRF key for custom implementation
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", os.Getenv("DB_DSN"), "Postgres connection string")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	//call the open db function
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	// Initialize the template cache
	templateCache, err := NewTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		addr:           addr,
		DB:             db, // <- Pass db here
		logger:         logger,
		MenuItem:       &data.MenuItemModel{DB: db},
		Order:          &data.OrderModel{DB: db},
		OrderItem:      &data.OrderItemModel{DB: db},
		Category:       &data.CategoryModel{DB: db},
		User:           &data.UserModel{DB: db},
		sessionStore:   sessionStore,
		Rating:         &data.RatingModel{DB: db},
		Analytics:      &data.AnalyticsModel{DB: db},
		Recommendation: &data.RecommendationModel{DB: db},
		templateCache:  templateCache,
		mailer: &mailer.Mailer{
			Host:     os.Getenv("MAIL_HOST"),
			Port:     587,
			Username: os.Getenv("MAIL_USERNAME"),
			Password: os.Getenv("MAIL_PASSWORD"),
			From:     os.Getenv("MAIL_USERNAME"),
		},
		env: os.Getenv("APP_ENV"),
	}
	app.logger.Info("application starting", "addr", *app.addr, "env", app.env)

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

// openDB function to open a connection to the database
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
