package main

import (
	"net/http"
	"os"

	"github.com/amilcar-vasquez/blessed-bites/backend/internal/handlers"
	"github.com/amilcar-vasquez/blessed-bites/backend/internal/middleware"
	"github.com/amilcar-vasquez/blessed-bites/backend/internal/realtime"
	"github.com/amilcar-vasquez/blessed-bites/backend/pkg/responses"
	"github.com/amilcar-vasquez/blessed-bites/internal/mailer"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	broker := realtime.NewBroker()

	menuHandler := handlers.NewMenuHandler(app.db)
	categoriesHandler := handlers.NewCategoriesHandler(app.db)
	resetMailer := &mailer.Mailer{
		Host:     os.Getenv("MAIL_HOST"),
		Port:     587,
		Username: os.Getenv("MAIL_USERNAME"),
		Password: os.Getenv("MAIL_PASSWORD"),
		From:     os.Getenv("MAIL_USERNAME"),
	}
	refreshSecret := []byte(app.config.jwtSecret)
	if v := os.Getenv("JWT_REFRESH_SECRET"); v != "" {
		refreshSecret = []byte(v)
	}
	cookieSecure := os.Getenv("APP_ENV") != "development"
	authHandler := handlers.NewAuthHandler(app.db, []byte(app.config.jwtSecret), refreshSecret, cookieSecure, resetMailer)
	ordersHandler := handlers.NewOrdersHandler(app.db, broker)
	streamHandler := handlers.NewStreamHandler(broker)
	ratingsHandler := handlers.NewRatingsHandler(app.db)
	adminHandler := handlers.NewAdminHandler(app.db)

	// Serve uploaded images from the monolith's static directory.
	// Strip the /uploads/ prefix and map to the configured uploads directory.
	uploadsDir := os.Getenv("UPLOADS_DIR")
	if uploadsDir == "" {
		uploadsDir = "uploads"
	}
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))

	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		responses.JSON(w, http.StatusOK, map[string]any{"status": "ok"})
	})

	mux.HandleFunc("GET /api/v1/menu", menuHandler.List)
	mux.HandleFunc("GET /api/v1/menu/{id}", menuHandler.GetByID)
	mux.HandleFunc("GET /api/v1/search", menuHandler.Search)
	mux.HandleFunc("GET /api/v1/categories", categoriesHandler.List)

	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/signup", authHandler.Signup)
	mux.HandleFunc("POST /api/v1/auth/logout", authHandler.Logout)
	mux.HandleFunc("POST /api/v1/auth/refresh", authHandler.Refresh)
	mux.HandleFunc("POST /api/v1/auth/reset-password-request", authHandler.ResetPasswordRequest)
	mux.HandleFunc("POST /api/v1/auth/reset-password", authHandler.ResetPassword)

	mux.Handle("POST /api/v1/orders", middleware.OptionalJWT([]byte(app.config.jwtSecret))(http.HandlerFunc(ordersHandler.Create)))
	mux.HandleFunc("GET /api/v1/orders/stream", streamHandler.OrdersStream)
	mux.Handle("POST /api/v1/ratings", middleware.OptionalJWT([]byte(app.config.jwtSecret))(http.HandlerFunc(ratingsHandler.Create)))
	mux.HandleFunc("GET /api/v1/ratings/{menu_item_id}", ratingsHandler.Average)

	adminProtected := func(h http.HandlerFunc) http.Handler {
		return middleware.RequireJWT([]byte(app.config.jwtSecret))(middleware.RequireAdmin(h))
	}
	mux.Handle("GET /api/v1/admin/menu", adminProtected(adminHandler.ListMenu))
	mux.Handle("POST /api/v1/admin/menu", adminProtected(adminHandler.CreateMenu))
	mux.Handle("PUT /api/v1/admin/menu/{id}", adminProtected(adminHandler.UpdateMenu))
	mux.Handle("DELETE /api/v1/admin/menu/{id}", adminProtected(adminHandler.DeleteMenu))
	mux.Handle("POST /api/v1/admin/category", adminProtected(adminHandler.CreateCategory))
	mux.Handle("DELETE /api/v1/admin/category/{id}", adminProtected(adminHandler.DeleteCategory))
	mux.Handle("GET /api/v1/admin/orders", adminProtected(adminHandler.ListOrders))
	mux.Handle("GET /api/v1/admin/users", adminProtected(adminHandler.ListUsers))
	mux.Handle("PUT /api/v1/admin/users/{id}", adminProtected(adminHandler.UpdateUser))
	mux.Handle("DELETE /api/v1/admin/users/{id}", adminProtected(adminHandler.DeleteUser))

	return middleware.CORS(app.config.corsOrigin)(middleware.Logging(app.logger)(mux))
}
