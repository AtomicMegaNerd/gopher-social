package main

import (
	"log"
	"log/slog"

	"github.com/atomicmeganerd/gopher-social/internal/db"
	"github.com/atomicmeganerd/gopher-social/internal/env"
	"github.com/atomicmeganerd/gopher-social/internal/store"
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
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "http://localhost:8080"),
		db: dbConfig{
			// postgres://user:password@host:port/dbname?sslmode=disable
			addr:         env.GetString("DATABASE_URL", ""), // no default, must be set
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 20),
			minIdleConns: env.GetInt("DB_MIN_IDLE_CONNS", 5),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env:     env.GetString("ENV", "development"),
		version: env.GetString("VERSION", "0.1.1"),
	}

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
	slog.Info("connected to database")
	store := store.NewPostgresStorage(pool)

	app := &application{config: cfg, store: store}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
