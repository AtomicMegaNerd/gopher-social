package main

import (
	"log/slog"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error(
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
	slog.Error(
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
	slog.Error(
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
