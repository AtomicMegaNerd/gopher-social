package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}

type PostStore struct {
	db *pgxpool.Pool
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	var createdAt, updatedAt time.Time

	if err := s.db.QueryRow(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		post.Tags, // pgx supports slices for array types directly
	).Scan(
		&post.ID,
		&createdAt,
		&updatedAt,
	); err != nil {
		return err
	}

	post.CreatedAt = createdAt.Format(time.RFC3339)
	post.UpdatedAt = updatedAt.Format(time.RFC3339)

	return nil
}

func (s *PostStore) GetByID(ctx context.Context, postID int64) (*Post, error) {
	query := `
		SELECT title, content, user_id, tags, created_at, updated_at
		FROM posts
		WHERE id=$1
	`

	var createdAt, updatedAt time.Time

	post := &Post{ID: postID}
	if err := s.db.QueryRow(ctx, query, postID).Scan(
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.Tags,
		&createdAt,
		&updatedAt,
	); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	post.CreatedAt = createdAt.Format(time.RFC3339)
	post.UpdatedAt = updatedAt.Format(time.RFC3339)

	return post, nil
}

func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id=$1`

	res, err := s.db.Exec(ctx, query, postID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET title=$1, content=$2, updated_at=$3
		WHERE id=$4
		RETURNING updated_at
	`

	var updatedAt time.Time
	if err := s.db.QueryRow(
		ctx,
		query,
		post.Title,
		post.Content,
		time.Now(),
		post.ID,
	).Scan(&updatedAt); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	post.UpdatedAt = updatedAt.Format(time.RFC3339)
	return nil
}
