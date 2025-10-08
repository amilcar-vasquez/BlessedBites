package main

import (
	"net/http"
)

// GET /checkout - render a dedicated checkout page showing the order sidebar
func (app *application) checkoutHandler(w http.ResponseWriter, r *http.Request) {
	data := app.addDefaultData(NewTemplateData(), w, r)

	// render checkout template
	if err := app.render(w, http.StatusOK, "checkout.tmpl", data); err != nil {
		app.logger.Error("Error rendering checkout template", "error", err)
	}
}
