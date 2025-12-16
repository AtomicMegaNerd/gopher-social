package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(
		"internal server error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	writeJSONError(
		w,
		http.StatusInternalServerError,
		"the server encountered a problem and could not process your request",
	)
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn(
		"bad request error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	writeJSONError(
		w,
		http.StatusBadRequest,
		err.Error(),
	)
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn(
		"not found error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	writeJSONError(
		w,
		http.StatusNotFound,
		"resource not found",
	)
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(
		"conflict error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error(),
	)
	writeJSONError(
		w,
		http.StatusConflict,
		"resource conflict",
	)
}

func (app *application) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn(
		"unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error(),
	)

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) unauthorizedBasicError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn(
		"unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error(),
	)

	// Set the WWW-Authenticate header to indicate that basic authentication is required
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}
