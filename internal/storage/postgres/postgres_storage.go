package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Pacahar/graphql-comments/internal/storage"
)

func NewPostgresStorage(dsn string) (*storage.Storage, error) {
	const op = "storage.postgres.NewPostgresStorage"

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	PostgresPostStorage, err := NewPostgresPostStorage(db)
	if err != nil {
		return nil, err
	}

	PostgresCommentStorage, err := NewPostgresCommentStorage(db)
	if err != nil {
		return nil, err
	}

	return &storage.Storage{
		Post:    PostgresPostStorage,
		Comment: PostgresCommentStorage,
	}, nil
}
