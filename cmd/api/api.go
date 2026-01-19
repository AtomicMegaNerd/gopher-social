package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/atomicmeganerd/gopher-social/docs"
	"github.com/atomicmeganerd/gopher-social/internal/auth"
	"github.com/atomicmeganerd/gopher-social/internal/mailer"
	"github.com/atomicmeganerd/gopher-social/internal/store"
	"github.com/atomicmeganerd/gopher-social/internal/store/cache"
)

const (
	httpTimeout  = 60 * time.Second
	writeTimeout = 30 * time.Second
	readTimeout  = 10 * time.Second
	idleTimeout  = 60 * time.Second
)

// The primary application struct
type application struct {
	config        config
	dbStore       *store.Storage
	cacheStore    *cache.Storage
	logger        *slog.Logger
	mailer        mailer.Client
	authenticator auth.Authenticator
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// TODO: Finishc configuring this properly
		AllowedOrigins:   []string{"http://localhost:5146"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

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

		// r.With applies the middleware to the route in a nice clean way
		r.With(app.BasicAuthMiddleware()).Get("/health", app.healthCheckHandler)

		// Swagger documentation route
		docsUrl := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsUrl)))

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Patch("/", app.checkPostOwnership("moderator", app.updatePostHandler))
				r.Delete("/", app.checkPostOwnership("admin", app.deletePostHandler))
				r.Post("/comments", app.createCommentHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		// routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {

	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
	}

	app.logger.Info("Starting server", "port", app.config.addr)
	return srv.ListenAndServe()
}
