package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/atomicmeganerd/gopher-social/internal/auth"
	"github.com/atomicmeganerd/gopher-social/internal/db"
	"github.com/atomicmeganerd/gopher-social/internal/env"
	"github.com/atomicmeganerd/gopher-social/internal/mailer"
	"github.com/atomicmeganerd/gopher-social/internal/store"
	"github.com/lmittmann/tint"
)

const version = "0.1.0"

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gohpers

//	@contact.name	AtomicMegaNerd
//	@contact.url	https://github.com/atomicmeganerd

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {

	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "http://localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			// postgres://user:password@host:port/dbname?sslmode=disable
			addr:         env.GetString("DATABASE_URL", ""), // no default, must be set
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 20),
			minIdleConns: env.GetInt("DB_MIN_IDLE_CONNS", 5),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
		version: env.GetString("VERSION", "0.1.1"),
		auth: authConfig{
			basic: basicAuthConfig{
				username: env.GetString("BASIC_USERNAME", ""),
				password: env.GetString("BASIC_PASSWORD", ""),
			},
			jwtToken: jwtTokenConfig{
				secret:    env.GetString("JWT_SECRET", ""),
				tokenHost: env.GetString("JWT_TOKEN_HOST", ""),
				expiry:    time.Hour * 24 * 3, // 3 days
			},
		},
	}

	// Init logging
	handler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
	logger := slog.New(handler)

	pool, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.minIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer pool.Close()
	logger.Info("connected to database")
	store := store.NewPostgresStorage(pool)

	mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.jwtToken.secret,
		cfg.auth.jwtToken.tokenHost,
		cfg.auth.jwtToken.tokenHost,
	)

	app := &application{
		config: cfg, store: store, mailer: mailer, logger: logger, authenticator: jwtAuthenticator,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
