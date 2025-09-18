package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Pacahar/graphql-comments/internal/models"
	storageErrors "github.com/Pacahar/graphql-comments/internal/storage/errors"
)

type PostPostgresStorage struct {
	db *sql.DB
}

func NewPostgresPostStorage(db *sql.DB) (*PostPostgresStorage, error) {
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

	return &PostPostgresStorage{db: db}, nil
}

func (ps *PostPostgresStorage) CreatePost(ctx context.Context, title, content string, commentsDisabled bool) (int64, error) {
	const op = "storage.postgres.post.CreatePost"

	var id int64
	err := ps.db.QueryRowContext(ctx, `
		INSERT INTO post (title, content, comments_disabled)
		VALUES ($1, $2, $3)
		RETURNING id`,
		title, content, commentsDisabled,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (ps *PostPostgresStorage) GetPostByID(ctx context.Context, id int64) (models.Post, error) {
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
			return models.Post{}, storageErrors.ErrPostNotFound
		}
		return models.Post{}, fmt.Errorf("%s: %w", op, err)
	}

	return post, nil
}

func (ps *PostPostgresStorage) GetAllPosts(ctx context.Context) ([]models.Post, error) {
	const op = "storage.postgres.post.GetAllPosts"

	posts := make([]models.Post, 0)

	rows, err := ps.db.QueryContext(ctx, `
		SELECT id, title, content, comments_disabled, created_at
		FROM post
		ORDER BY created_at ASC`,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.CommentsDisabled,
			&post.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: iteration failed: %w", op, err)
	}

	return posts, nil
}

func (ps *PostPostgresStorage) DeletePost(ctx context.Context, id int64) error {
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
