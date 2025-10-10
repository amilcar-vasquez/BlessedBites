// file: cmd/web/render.go
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *TemplateData) error {

	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.logger.Error("Template does not exist", "template", page, "error", err)

		return err
	}

	err := ts.Execute(buf, data)
	if err != nil {
		err := fmt.Errorf("failed to render template: %s", err)
		app.logger.Error("failed to render template", "template", page, "error", err)

		return err
	}

	// Set the status header before writing the body (correct order).
	w.WriteHeader(status)

	// Write the rendered template to the ResponseWriter.
	_, err = buf.WriteTo(w)
	if err != nil {
		// Log write failures at Info to avoid causing handlers to attempt
		// secondary writes (which would result in superfluous WriteHeader
		// logs). In practice these are often client disconnects.
		app.logger.Info("failed to write template to browser", "error", err)
		return nil
	}

	return nil
}
