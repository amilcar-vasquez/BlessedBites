package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type config struct {
	addr       string
	dsn        string
	jwtSecret  string
	corsOrigin string
}

type application struct {
	config config
	logger *slog.Logger
	db     *sql.DB
}

func main() {
	cfg := config{}
	flag.StringVar(&cfg.addr, "addr", ":8080", "API listen address")
	flag.StringVar(&cfg.dsn, "dsn", os.Getenv("DB_DSN"), "Postgres DSN")
	flag.StringVar(&cfg.jwtSecret, "jwt-secret", os.Getenv("JWT_SECRET"), "JWT HMAC secret")
	flag.StringVar(&cfg.corsOrigin, "cors-origin", getenvDefault("CORS_ORIGIN", "*"), "Allowed CORS origin")
	flag.Parse()

	if cfg.dsn == "" {
		fmt.Fprintln(os.Stderr, "missing DB_DSN / -dsn")
		os.Exit(1)
	}
	if cfg.jwtSecret == "" {
		fmt.Fprintln(os.Stderr, "missing JWT_SECRET / -jwt-secret")
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("starting api", "addr", cfg.addr)

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	app := &application{config: cfg, logger: logger, db: db}
	if err := app.serve(); err != nil {
		logger.Error("api server failed", "error", err)
		os.Exit(1)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func getenvDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
