package store

import (
	"context"
	"database/sql"
)

type users struct {
	ID        int8   `json:"id"`
	Username  string `json:"username"`
	Age       int8   `json:"age"`
	Password  string `json:"_"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
type userStore struct {
	db *sql.DB
}

func (s *userStore) Create(ctx context.Context, users *users) error {
	query := `
	INSERT INTO user (user.username,user.password,user.age)
	VALUES ($1, $2, $3,$4)
	RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRowContext(ctx, query, users.Username, users.Age, users.Password, users.Email).Scan(
		&users.CreatedAt,
		&users.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}
