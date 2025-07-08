package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

/*
	 Ported from this original code as the pq libarary is now deprecated:

	 package db

	 import (
		"context"
		"database/sql"
		"time"
	 )

	 func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
		db, err := sql.Open("postgres", addr)
		if err != nil {
			return nil, err
		}
		db.SetMaxOpenConns(maxOpenConns)
		db.SetMaxIdleConns(maxIdleConns)
		duration, err := time.ParseDuration(maxIdleTime)
		if err != nil {
			return nil, err
		}
		db.SetConnMaxIdleTime(duration)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			return nil, err
		}

		return db, nil
	 }
*/
func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(addr)
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(maxOpenConns)
	config.MinIdleConns = int32(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}

	config.MaxConnIdleTime = duration

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
