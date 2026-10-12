package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

// Post represents a post in the social media application.
type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostRepository interface {
	Create(context.Context, *Post) error
}

type postStore struct {
	db *sql.DB
}

func (s *postStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (title, content, user_id, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, status, created_at, updated_at
	`

	return s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.Status,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
}
