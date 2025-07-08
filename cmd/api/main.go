package main

import (
	"log"

	"github.com/atomicmeganerd/rcd-gopher-social/internal/db"
	"github.com/atomicmeganerd/rcd-gopher-social/internal/env"
	"github.com/atomicmeganerd/rcd-gopher-social/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			// postgres://user:password@host:port/dbname?sslmode=disable
			addr:         env.GetString("DB_ADDR", ""), // no default, must be set
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15min"),
		},
	}

	pool, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	store := store.NewPostgresStorage(pool)

	app := &application{config: cfg, store: store}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
