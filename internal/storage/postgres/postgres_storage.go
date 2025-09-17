package postgres

import (
	"database/sql"
	"fmt"
)

type PostgresStorage struct {
	postStorage    *PostPostgresStorage
	commentStorage *CommentPostgresStorage
}

func NewPostgresStorage(storagePath string) (*PostgresStorage, error) {
	const op = "storage.postgres.NewPostgresStorage"

	db, err := sql.Open("postgres", storagePath)

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

	return &PostgresStorage{
		postStorage:    PostgresPostStorage,
		commentStorage: PostgresCommentStorage,
	}, nil
}
