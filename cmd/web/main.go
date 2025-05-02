// filepath: cmd/web/main.go
package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"html/template"
	"log/slog"
	"os"
	"time"
)

var sessionStore = sessions.NewCookieStore([]byte("super-secret-key"))

type application struct {
	addr            *string
	MenuItem        *data.MenuItemModel
	Order           *data.OrderModel
	Category        *data.CategoryModel
	User            *data.UserModel
	sessionStore    *sessions.CookieStore
	Analytics       *data.AnalyticsModel
	Recommendations *data.RecommendationModel
	logger          *slog.Logger
	templateCache   map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "postgres://blessed_bites:Matthew.5:6@localhost/blessed_bites?sslmode=disable", "Postgres connection string")

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
		addr:            addr,
		logger:          logger,
		MenuItem:        &data.MenuItemModel{DB: db},
		Order:           &data.OrderModel{DB: db},
		Category:        &data.CategoryModel{DB: db},
		User:            &data.UserModel{DB: db},
		sessionStore:    sessionStore,
		Analytics:       &data.AnalyticsModel{DB: db},
		Recommendations: &data.RecommendationModel{DB: db},
		templateCache:   templateCache,
	}

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
