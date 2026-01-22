package cache

import (
	"context"

	"github.com/atomicmeganerd/gopher-social/internal/store"
)

func NewMockStore() *Storage {
	return &Storage{
		Users: &MockUsersCacheStorage{},
	}
}

type MockUsersCacheStorage struct{}

func (m *MockUsersCacheStorage) Get(ctx context.Context, id int64) (*store.User, error) {
	return nil, nil
}

func (m *MockUsersCacheStorage) Set(ctx context.Context, user *store.User) error {
	return nil
}
