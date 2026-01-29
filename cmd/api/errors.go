package main

import (
	"fmt"
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
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/WWW-Authenticate
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) forbiddenError(w http.ResponseWriter, r *http.Request) {
	app.logger.Warn("forbidden error", "method", r.Method, "path", r.URL.Path)

	writeJSONError(w, http.StatusForbidden, "forbidden")
}

func (app *application) rateLimitExceeededError(
	w http.ResponseWriter, r *http.Request, retryAfter string,
) {

	app.logger.Warn(
		"rate limited exceeded error",
		"method", r.Method,
		"path", r.URL.Path,
		"sourceIP", r.RemoteAddr,
	)

	w.Header().Set("Retry-After", retryAfter)
	writeJSONError(
		w, http.StatusTooManyRequests, fmt.Sprintf("rate limit exceeded, retry after: %s", retryAfter),
	)
}
