package cache

import (
	"context"

	"github.com/atomicmeganerd/gopher-social/internal/store"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
	}
}

func NewCacheStorage(rds *redis.Client) *Storage {
	return &Storage{
		Users: &UserStore{
			rds: rds,
		},
	}
}
