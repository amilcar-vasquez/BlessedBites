package main

import (
	"context"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"os"
	"time"

	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	_ "github.com/lib/pq"
)

type application struct {
	addr            *string
	MenuItems       *data.MenuItemModel
	Orders          *data.OrderModel
	Categories      *data.CategoryModel
	Users           *data.UserModel
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

	templateCache, err := NewTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		addr:   addr,
		logger: logger,

		Orders:          &data.OrderModel{DB: db},
		Categories:      &data.CategoryModel{DB: db},
		Users:           &data.UserModel{DB: db},
		Analytics:       &data.AnalyticsModel{DB: db},
		Recommendations: &data.RecommendationModel{DB: db},
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
