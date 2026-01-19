package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/atomicmeganerd/gopher-social/internal/store"
	"github.com/redis/go-redis/v9"
)

type UserStore struct {
	rds *redis.Client
}

const UserExpTime = time.Minute

func (u *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	// If cache is not enabled
	if u.rds == nil {
		return nil, nil
	}

	cacheKey := fmt.Sprintf("user-%v", userID)

	data, err := u.rds.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (u *UserStore) Set(ctx context.Context, user *store.User) error {

	if user == nil {
		return errors.New("user cannot be nil")
	}

	if u.rds == nil {
		return nil
	}

	cacheKey := fmt.Sprintf("user-%v", user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = u.rds.Set(ctx, cacheKey, json, UserExpTime).Result()
	if err != nil {
		return err
	}

	return nil
}
