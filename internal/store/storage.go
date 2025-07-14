package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
	}
	Users interface {
		Create(context.Context, *User) error
	}
}

func NewPostgresStorage(db *pgxpool.Pool) *Storage {
	return &Storage{
		Posts: &PostsStore{db},
		Users: &UsersStore{db},
	}
}
