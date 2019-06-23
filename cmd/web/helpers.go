package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	app.writeJSON(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	app.writeJSON(w, status, http.StatusText(status))
}

func (app *application) writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}