package cache

import (
	"context"

	"github.com/atomicmeganerd/gopher-social/internal/store"
	"github.com/stretchr/testify/mock"
)

func NewMockStore() *Storage {
	return &Storage{
		Users: &MockUsersCacheStorage{},
	}
}

type MockUsersCacheStorage struct {
	mock.Mock
}

func (m *MockUsersCacheStorage) Get(ctx context.Context, userID int64) (*store.User, error) {
	args := m.Called(userID)
	return nil, args.Error(1)
}

func (m *MockUsersCacheStorage) Set(ctx context.Context, user *store.User) error {
	args := m.Called(user)
	return args.Error(1)
}
