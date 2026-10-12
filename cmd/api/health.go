package main

import (
	"net/http"
)

// healthHandler returns a simple JSON health check response
func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	env:=app.config.env
	if env==""{
		env="development"
	}
	data := map[string]string{
		"status":  "ok",
		"env":     env,
		"version": version,
	}

	if err := writeJSON(w, http.StatusOK, data); err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
