package main

import (
	"context"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"os"
	"time"

	"github.com/amilcar-vasquez/blessed-bites"
)

type application struct {
	address      		*string
	Orders         		*data.OrderModel
    Categories     		*data.CategoryModel
    Users          		*data.UserModel
    Analytics      		*data.AnalyticsModel
    Recommendations 	*data.RecommendationModel
	logger 	 			*slog.Logger
	templateCache 		map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "postgres://blessed:tapir@localhost/journal?sslmode=disable", "Postgres connection string")