package store

import (
	"context"
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
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

// This will be used within the feed
type PostWithMetadata struct {
	Post             // Composition is not inheritance :-)
	CommentCount int `json:"comments_count"`
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
		SELECT title, content, user_id, tags, created_at, updated_at, version
		FROM posts
		WHERE id=$1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var createdAt, updatedAt time.Time

	post := &Post{ID: postID}
	if err := s.db.QueryRow(ctx, query, postID).Scan(
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.Tags,
		&createdAt,
		&updatedAt,
		&post.Version,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
	// Optimistic locking: only update if the version matches
	// This prevents lost updates in concurrent scenarios
	// The version is incremented on each successful update
	query := `
		UPDATE posts
		SET title=$1, content=$2, updated_at=$3, version=version + 1
		WHERE id=$4 AND version=$5
		RETURNING updated_at, version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var updatedAt time.Time
	if err := s.db.QueryRow(
		ctx,
		query,
		post.Title,
		post.Content,
		time.Now(),
		post.ID,
		post.Version,
	).Scan(&updatedAt, &post.Version); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return ErrNotFound
		default:
			return err
		}
	}

	post.UpdatedAt = updatedAt.Format(time.RFC3339)
	return nil
}

func (s *PostStore) GetUserFeed(
	ctx context.Context, userID int64, pq PaginatedFeedQuery,
) ([]PostWithMetadata, error) {

	query := `
		SELECT
			p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
	    	u.username,
			COUNT(c.id) AS comment_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
	  	LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
	  	WHERE f.user_id = $1 OR p.user_id = $1
	   	GROUP BY p.id, u.username
		ORDER BY p.created_at ` + pq.Sort + `
		LIMIT $2 OFFSET $3
		;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.Query(
		ctx,
		query,
		userID,
		pq.Limit,
		pq.Offset,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostWithMetadata
	var createdAt time.Time
	for rows.Next() {
		var p PostWithMetadata
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&createdAt,
			&p.Version,
			&p.Tags,
			&p.User.Username,
			&p.CommentCount,
		)
		if err != nil {
			return nil, err
		}

		p.CreatedAt = createdAt.Format(time.RFC3339)

		feed = append(feed, p)
	}

	return feed, nil
}
