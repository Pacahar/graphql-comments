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

type CommentPostgresStorage struct {
	db *sql.DB
}

func NewPostgresCommentStorage(db *sql.DB) (*CommentPostgresStorage, error) {
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

	return &CommentPostgresStorage{db: db}, nil
}

func (cs *CommentPostgresStorage) CreateComment(ctx context.Context, content string, postID int64, parentID *int64) (int64, error) {
	const op = "storage.postgres.comment.CreateComment"

	var id int64
	err := cs.db.QueryRowContext(ctx, `
		INSERT INTO comment (content, post_id, parent_id)
		VALUES ($1, $2, $3)
		RETURNING id`,
		content, postID, parentID,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (cs *CommentPostgresStorage) GetCommentByID(ctx context.Context, id int64) (models.Comment, error) {
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
			return models.Comment{}, storageErrors.ErrCommentNotFound
		}
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}

func (cs *CommentPostgresStorage) GetCommentsByParentID(ctx context.Context, ParentID int64) ([]models.Comment, error) {
	const op = "storage.postgres.comment.GetCommentsByParentID"

	rows, err := cs.db.QueryContext(ctx, `
		SELECT id, post_id, parent_id, content, created_at
		FROM comment
		WHERE parent_id = $1
		ORDER BY created_at ASC`,
		ParentID,
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

func (cs *CommentPostgresStorage) GetCommentsByPostID(ctx context.Context, postID int64, limit *int64, offset *int64) ([]models.Comment, error) {
	const op = "storage.postgres.comment.GetCommentsByPostID"

	var rows *sql.Rows
	var err error

	if limit != nil && offset != nil {
		rows, err = cs.db.QueryContext(ctx, `
		SELECT id, post_id, parent_id, content, created_at
		FROM comment
		WHERE post_id = $1
		AND parent_id IS NULL
		ORDER BY created_at ASC
		LIMIT $2
		OFFSET $3
	`, postID, *limit, *offset)
	} else {
		rows, err = cs.db.QueryContext(ctx, `
		SELECT id, post_id, parent_id, content, created_at
		FROM comment
		WHERE post_id = $1
		AND parent_id IS NULL
		ORDER BY created_at ASC`, postID)
	}

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

func (cs *CommentPostgresStorage) DeleteComment(ctx context.Context, id int64) error {
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

func (cs *CommentPostgresStorage) DeleteCommentsByPostID(ctx context.Context, postID int64) error {
	const op = "storage.postgres.comment.DeleteCommentsByPostID"

	_, err := cs.db.ExecContext(ctx, `
		DELETE FROM comment
		WHERE post_id=$1`,
		postID,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
