package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/atomicmeganerd/gopher-social/docs"
	"github.com/atomicmeganerd/gopher-social/internal/auth"
	"github.com/atomicmeganerd/gopher-social/internal/mailer"
	"github.com/atomicmeganerd/gopher-social/internal/ratelimiter"
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
	rateLimiter   ratelimiter.Limiter
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	// This middleware recovers from panics and writes a 500 if there is one.
	r.Use(middleware.Recoverer)
	// NOTE: This middleware logs the IP address of the requestor. This is crucial for our
	// rate limiter
	r.Use(middleware.RealIP)
	// This middleware adds a request ID to each request.
	r.Use(middleware.RequestID)
	// What a great way to set timeout!
	r.Use(middleware.Timeout(httpTimeout))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{app.config.frontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Our custom middleware, make sure rate limter comes last
	r.Use(app.RateLimiterMiddleware)

	// Creating the routes is really easy with chi.
	r.Route("/v1", func(r chi.Router) {

		// Do not use basic auth anymore due to need for graceful shutdown
		r.Get("/health", app.healthCheckHandler)

		// This is provided
		r.With(app.BasicAuthMiddleware()).Get("/debug/vars", expvar.Handler().ServeHTTP)

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

	shutdown := make(chan error)

	go func() {

		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Warn("Signal recieved, initiating shutdown", "signal", s.String())
		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Info(
		"Starting GopherSocial server instance",
		"port", app.config.addr,
		"env", app.config.env,
	)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		app.logger.Error("Unexpected error...", "error", err)
		return err
	}

	err = <-shutdown
	if err != nil {
		app.logger.Error("Unexpected error...", "error", err)
		return err
	}

	app.logger.Info(
		"server has stopped with no errors", "addr", app.config.addr, "env", app.config.env,
	)
	return nil
}
