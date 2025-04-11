// file: cmd/web/routes.go
package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	//create handlers for all the routes
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /signup", app.signupForm)
	mux.HandleFunc("POST /signup/new", app.signupHandler)
	mux.HandleFunc("GET /menu/add", app.addMenuItemForm)
	mux.HandleFunc("POST /menu/add/new", app.addMenuItemHandler)
	mux.HandleFunc("GET /menu", app.menuPageHandler)
	mux.HandleFunc("POST /menu/delete", app.deleteMenuItem)
	mux.HandleFunc("POST /menu/edit", app.editMenuItem)
	mux.HandleFunc("POST /menu/update", app.updateMenuItem)
	mux.HandleFunc("POST /category/add", app.addCategory)
	mux.HandleFunc("POST /category/delete", app.deleteCategory)

	return app.loggingMiddleware(mux)
}
