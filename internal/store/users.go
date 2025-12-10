package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash

	return nil
}

type UserStore struct {
	db *pgxpool.Pool
}

func (s *UserStore) Create(ctx context.Context, tx pgx.Tx, user *User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var createdAt time.Time
	if err := tx.QueryRow(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password.hash,
	).Scan(
		&user.ID,
		&createdAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				switch pgErr.ConstraintName {
				case "users_email_key":
					return ErrDuplicateEmail
				case "users_username_key":
					return ErrDuplicateUsername
				}
			}
		}
	}

	user.CreatedAt = createdAt.Format(time.RFC3339)
	return nil
}

func (s *UserStore) GetByID(ctx context.Context, id string) (*User, error) {
	// WARNING: Do NOT include the password field in the SELECT statement
	query := `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	var createdAt time.Time
	if err := s.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&createdAt,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	user.CreatedAt = createdAt.Format(time.RFC3339)
	return user, nil
}

func (s *UserStore) CreateAndInvite(
	ctx context.Context, user *User, token string, inviteExpiry time.Duration,
) error {
	return withTx(s.db, ctx, func(tx pgx.Tx) error {
		// create the user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		// create the user invitation
		if err := s.createUserInvitation(ctx, tx, token, inviteExpiry, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) createUserInvitation(
	ctx context.Context, tx pgx.Tx, token string, inviteExpiry time.Duration, userID int64,
) error {

	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.Exec(ctx, query, token, userID, time.Now().Add(inviteExpiry))
	if err != nil {
		return err
	}

	return nil
}
