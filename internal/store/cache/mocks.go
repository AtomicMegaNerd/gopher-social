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
	// We set this to 0 because the user argument will be nil here in the test
	return args.Error(0)
}
