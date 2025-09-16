package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	storageErrors "github.com/Pacahar/graphql-comments/internal/storage/errors"
	_ "github.com/lib/pq"

	"github.com/Pacahar/graphql-comments/internal/models"
)

type PostgresCommentStorage struct {
	db *sql.DB
}

func NewPostgresCommentStorage(db *sql.DB) (*PostgresCommentStorage, error) {
	const op = "storage.postgres.NewPostgresCommentStorage"

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS comment(
			id SERIAL PRIMARY KEY,
			post_id INTEGER NOT NULL,
			parent_id INTEGER NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW() NOT NULL,
			FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE,
			FOREIGN KEY (parent_id) REFERENCES comment(id) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_comment_post_id ON comment(post_id);
		CREATE INDEX IF NOT EXISTS idx_comment_parent_id ON comment(parent_id);
		CREATE INDEX IF NOT EXISTS idx_comment_created_at ON comment(created_at);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &PostgresCommentStorage{db: db}, nil
}

func (cs *PostgresCommentStorage) CreateComment(ctx context.Context, content string, postID int64, parentID *int64) error {
	const op = "storage.postgres.comment.CreateComment"

	var err error

	_, err = cs.db.ExecContext(ctx, `
		INSERT INTO comment(content, post_id, parent_id)
		VALUES($1, $2, $3)
	`, content, postID, *parentID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (cs *PostgresCommentStorage) GetCommentByID(ctx context.Context, id int64) (*models.Comment, error) {
	const op = "storage.postgres.comment.GetCommentByID"

	comment := models.Comment{}

	row := cs.db.QueryRowContext(ctx, `
		SELECT id, post_id, parent_id, content, created_at 
		FROM comment 
		WHERE id=$1`,
		id,
	)

	err := row.Scan(
		&comment.ID,
		&comment.PostID,
		&comment.ParentID,
		&comment.Content,
		&comment.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storageErrors.ErrCommentNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &comment, nil
}

func (cs *PostgresCommentStorage) GetCommentsByPostID(ctx context.Context, postID int64) ([]models.Comment, error) {
	const op = "storage.postgres.comment.GetCommentsByPostID"

	rows, err := cs.db.QueryContext(ctx, `
		SELECT id, post_id, parent_id, content, created_at
		FROM comment
		WHERE post_id=$1
		ORDER BY created_at ASC`,
		postID,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	comments := make([]models.Comment, 0)

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.ParentID,
			&comment.Content,
			&comment.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: iteration failed: %w", op, err)
	}

	return comments, nil
}

func (cs *PostgresCommentStorage) DeleteComment(ctx context.Context, id int64) error {
	const op = "storage.postgres.comment.DeleteComment"

	_, err := cs.db.ExecContext(ctx, `
		DELETE FROM comment
		WHERE id=$1`,
		id,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
