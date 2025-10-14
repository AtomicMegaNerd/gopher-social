package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"` // Password should not be exposed in JSON
	CreatedAt string `json:"created_at"`
}

type UserStore struct {
	db *pgxpool.Pool
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	if err := s.db.QueryRow(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	); err != nil {
		return err
	}
	return nil
}
