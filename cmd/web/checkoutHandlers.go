package main

import (
	"net/http"
)

// GET /checkout - render a dedicated checkout page showing the order sidebar
func (app *application) checkoutHandler(w http.ResponseWriter, r *http.Request) {
	data := app.addDefaultData(NewTemplateData(), w, r)

	// Read flash session and clear it so client can respond (e.g., clear localStorage)
	if flashSession, err := app.getSessionSafe(w, r, "flash"); err == nil && flashSession != nil {
		if msg, ok := flashSession.Values["alertMessage"].(string); ok && msg != "" {
			data.AlertMessage = msg
			if typ, ok := flashSession.Values["alertType"].(string); ok && typ != "" {
				data.AlertType = typ
			} else {
				data.AlertType = "alert-info"
			}
			flashSession.Options.MaxAge = -1
			_ = flashSession.Save(r, w)
		}
	}

	// render checkout template
	if err := app.render(w, http.StatusOK, "checkout.tmpl", data); err != nil {
		app.logger.Error("Error rendering checkout template", "error", err)
	}
}
