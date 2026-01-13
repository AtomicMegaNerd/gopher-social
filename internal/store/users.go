package store

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	IsActive  bool     `json:"is_active"`
}

type password struct {
	text *string
	hash []byte
}

type UserInvitation struct {
	Token  string `json:"token"`
	UserID int64  `json:"user_id"`
	Expiry string `json:"expiry"`
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
	query := /* sql */ `
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

func (s *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	// WARNING: Do NOT include the password field in the SELECT statement
	query := /* sql */ `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = $1
		AND is_active = true
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

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := /* sql */ `
		SELECT id, username, email, password, created_at
		FROM users
		WHERE email = $1
		AND is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	var createdAt time.Time
	if err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
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

func (s *UserStore) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx pgx.Tx) error {

		// Get the user
		user, err := s.getUserFromToken(ctx, tx, token)
		if err != nil {
			return err
		}

		// Update the user
		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}

		// Delete all pending invitations for that user
		if err := s.deleteUserInvitations(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) Delete(ctx context.Context, userID int64) error {

	return withTx(s.db, ctx, func(tx pgx.Tx) error {
		if err := s.delete(ctx, tx, userID); err != nil {
			return err
		}

		if err := s.deleteUserInvitations(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) update(ctx context.Context, tx pgx.Tx, user *User) error {

	query := /* sql */ `
	  UPDATE users u
	  SET username=$1, email=$2, is_active=$3
	  WHERE u.id = $4
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.Exec(
		ctx, query, user.Username, user.Email, user.IsActive, user.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) delete(ctx context.Context, tx pgx.Tx, userID int64) error {

	query := /* sql */ `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) createUserInvitation(
	ctx context.Context, tx pgx.Tx, token string, inviteExpiry time.Duration, userID int64,
) error {

	query := /* sql */ `INSERT INTO user_invitations
							(token, user_id, expiry)
							VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.Exec(ctx, query, token, userID, time.Now().Add(inviteExpiry))
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) deleteUserInvitations(ctx context.Context, tx pgx.Tx, userID int64) error {
	query := /* sql */ `DELETE FROM user_invitations WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) getUserFromToken(
	ctx context.Context, tx pgx.Tx, token string,
) (*User, error) {

	query := /* sql */ `SELECT u.id, u.username, u.email, u.created_at, u.is_active FROM users u
							JOIN user_invitations i
							ON u.id = i.user_id
							WHERE i.token = $1 and i.expiry > $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	var createdAt time.Time
	user := &User{}
	if err := tx.QueryRow(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&createdAt,
		&user.IsActive,
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
