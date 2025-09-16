package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Pacahar/graphql-comments/internal/models"
	storageErrors "github.com/Pacahar/graphql-comments/internal/storage/errors"
)

type PostgresPostStorage struct {
	db *sql.DB
}

func NewPostgresPostStorage(db *sql.DB) (*PostgresPostStorage, error) {
	const op = "storage.postgres.NewPostgresPostStorage"

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS post(
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			comments_disabled BOOLEAN NOT NULL,
			created_at TIMESTAMP DEFAULT NOW() NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_post_created_at ON post(created_at);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &PostgresPostStorage{db: db}, nil
}

func (ps *PostgresPostStorage) CreatePost(ctx context.Context, title, content string, comments_disabled bool) error {
	const op = "storage.postgres.post.CreatePost"

	_, err := ps.db.ExecContext(ctx, `
		INSERT INTO post(title, content, comments_disabled)
		VALUES($1, $2, $3)`,
		title,
		content,
		comments_disabled,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (ps *PostgresPostStorage) GetPostByID(ctx context.Context, id int64) (*models.Post, error) {
	const op = "storage.postgres.post.GetPostByID"

	post := models.Post{}

	row := ps.db.QueryRowContext(ctx, `
		SELECT id, title, content, comments_disabled, created_at 
		FROM post 
		WHERE id=$1`,
		id,
	)

	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.CommentsDisabled,
		&post.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storageErrors.ErrPostNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &post, nil
}

func (ps *PostgresPostStorage) DeletePost(ctx context.Context, id int64) error {
	const op = "storage.postgres.post.DeletePost"

	_, err := ps.db.ExecContext(ctx, `
		DELETE FROM post
		WHERE id=$1`,
		id,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
