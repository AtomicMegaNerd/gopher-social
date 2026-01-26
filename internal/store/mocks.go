package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func NewMockStore() *Storage {
	return &Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
}

func (m *MockUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	return &User{
		ID: 42,
	}, nil
}

func (m *MockUserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	return nil, nil
}

func (m *MockUserStore) Create(ctx context.Context, tx pgx.Tx, user *User) error {
	return nil
}

func (m *MockUserStore) CreateAndInvite(
	ctx context.Context, user *User, token string, expiryDuration time.Duration,
) error {
	return nil
}

func (m *MockUserStore) Activate(ctx context.Context, token string) error {
	return nil
}

func (m *MockUserStore) Delete(ctx context.Context, id int64) error {
	return nil
}
