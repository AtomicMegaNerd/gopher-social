package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/atomicmeganerd/gopher-social/internal/auth"
	"github.com/atomicmeganerd/gopher-social/internal/db"
	"github.com/atomicmeganerd/gopher-social/internal/mailer"
	"github.com/atomicmeganerd/gopher-social/internal/store"
	"github.com/atomicmeganerd/gopher-social/internal/store/cache"
	"github.com/lmittmann/tint"
	"github.com/redis/go-redis/v9"
)

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

	cfg := NewConfig()

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
	dbStore := store.NewPostgresStorage(pool)

	var rds *redis.Client
	if cfg.cache.enabled {
		rds = cache.NewRedisClient(cfg.cache.addr, cfg.cache.password, cfg.cache.db)
		logger.Info("connected to redis cache")
	}
	cacheStore := cache.NewCacheStorage(rds)

	mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.jwtToken.secret,
		cfg.auth.jwtToken.tokenHost,
		cfg.auth.jwtToken.tokenHost,
	)

	app := &application{
		config:        cfg,
		cacheStore:    cacheStore,
		dbStore:       dbStore,
		mailer:        mailer,
		logger:        logger,
		authenticator: jwtAuthenticator,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
