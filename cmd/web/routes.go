// file: cmd/web/routes.go
package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func (app *application) mountStatic(mux *http.ServeMux, routePrefix, dirPath string) {
	wd, err := os.Getwd()
	if err != nil {
		app.logger.Error("failed to get working directory", "error", err)
		return
	}
	staticDir := filepath.Join(wd, dirPath)
	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle(routePrefix, http.StripPrefix(routePrefix, noCacheMiddleware(fileServer)))
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	app.mountStatic(mux, "/static/", "ui/static")
	app.mountStatic(mux, "/ui/static/", "ui/static")

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /signup", app.signupForm)
	mux.HandleFunc("POST /signup/new", app.signupHandler)
	mux.HandleFunc("GET /signup-thanks", app.signupThanks)
	mux.HandleFunc("GET /users", app.requireAuth(app.userPageHandler))            // list all users
	mux.HandleFunc("POST /user/update/form", app.requireAuth(app.updateUserForm)) //display the same signup form but with update values
	mux.HandleFunc("POST /user/update", app.requireAuth(app.updateUser))
	mux.HandleFunc("POST /users/delete", app.requireAuth(app.deleteUser))
	mux.HandleFunc("GET /login", app.loginForm)
	mux.HandleFunc("POST /login", app.loginHandler)
	mux.HandleFunc("POST /logout", app.logoutHandler)
	mux.HandleFunc("GET /menu/add", app.requireAuth(app.addMenuItemForm))
	mux.HandleFunc("POST /menu/add/new", app.requireAuth(app.addMenuItemHandler))
	mux.HandleFunc("GET /menu", app.requireAuth(app.menuPageHandler))
	mux.HandleFunc("POST /menu/delete", app.requireAuth(app.deleteMenuItem))
	mux.HandleFunc("POST /menu/edit", app.requireAuth(app.editMenuItem))
	mux.HandleFunc("POST /menu/update", app.requireAuth(app.updateMenuItem))
	mux.HandleFunc("POST /category/add", app.addCategory)
	mux.HandleFunc("POST /category/delete", app.deleteCategory)
	mux.HandleFunc("/search", app.searchMenuHandler)
	mux.HandleFunc("/search.json", app.searchMenuJSONHandler)
	mux.HandleFunc("/menu/category/", app.viewMenuByCategory)
	mux.HandleFunc("POST /orders", app.createOrderHandler)
	mux.HandleFunc("POST /ratings", app.createRatingHandler)
	mux.HandleFunc("GET /ratings/{menu_item_id}", app.getAverageRatingHandler)
	mux.HandleFunc("POST /menu/active", app.requireAuth(app.setActiveStatus))

	return app.loggingMiddleware(mux)
}
