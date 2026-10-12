package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Post interface {
		Create(context.Context, *Post) error
	}
	User interface {
		Create(context.Context, *users) error
	}
}

func NewSQLStorage(db *sql.DB) Storage {
	return Storage{
		Post: &postStore{db: db},
		User: &userStore{db: db},
	}
}