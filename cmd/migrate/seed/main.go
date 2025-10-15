package main

import (
	"log"
	"log/slog"

	"github.com/atomicmeganerd/rcd-gopher-social/internal/db"
	"github.com/atomicmeganerd/rcd-gopher-social/internal/env"
	"github.com/atomicmeganerd/rcd-gopher-social/internal/store"
)

func main() {
	addr := env.GetString("DATABASE_URL", "") // no default, must be set
	pool, err := db.New(
		addr,
		30,
		10,
		"15m",
	)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer pool.Close()
	slog.Info("connected to database")
	store := store.NewPostgresStorage(pool)

	db.Seed(store)
}
