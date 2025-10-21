package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/atomicmeganerd/gopher-social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type config struct {
	addr    string
	db      dbConfig
	env     string
	version string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	minIdleConns int
	maxIdleTime  string
}

const (
	httpTimeout  = 60 * time.Second
	writeTimeout = 30 * time.Second
	readTimeout  = 10 * time.Second
	idleTimeout  = 60 * time.Second
)

type application struct {
	config config
	store  *store.Storage
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	// This middleware recovers from panics and writes a 500 if there is one.
	r.Use(middleware.Recoverer)
	// This middleware logs the IP address of the requestor.
	r.Use(middleware.RealIP)
	// This middleware adds a request ID to each request.
	r.Use(middleware.RequestID)

	// What a great way to set timeout!
	r.Use(middleware.Timeout(httpTimeout))

	// Creating the routes is really easy with chi.
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
				r.Post("/comments", app.createCommentHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)
				r.Get("/", app.getUserHandler)

				// follow and unfollow routes
				// PUT /v1/users/{userID}/follow
				// PUT /v1/users/{userID}/unfollow
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
	}

	slog.Info("Starting server", "port", app.config.addr)
	return srv.ListenAndServe()
}
