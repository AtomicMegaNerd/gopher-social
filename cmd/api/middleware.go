package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedError(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedError(w, r, fmt.Errorf("invalid authorization header format"))
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedError(w, r, err)
				return
			}

			username := app.config.auth.basic.username
			password := app.config.auth.basic.password

			if username == "" || password == "" {
				app.unauthorizedError(w, r, fmt.Errorf("basic auth credentials are not set"))
			}

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.unauthorizedError(w, r, fmt.Errorf("invalid credentials"))
			}

			next.ServeHTTP(w, r)
		})
	}
}
