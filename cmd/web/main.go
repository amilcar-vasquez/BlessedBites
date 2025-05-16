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
