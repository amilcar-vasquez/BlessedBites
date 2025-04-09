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

	return app.loggingMiddleware(mux)
}
