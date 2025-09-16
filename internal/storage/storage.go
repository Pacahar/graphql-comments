package storage

import (
	"context"

	"github.com/Pacahar/graphql-comments/internal/models"
)

const (
	memoryStorage   = "memory"
	postgresStorage = "postgres"
)

type Storage struct {
	Comment CommentStorage
	Post    PostStorage
}

type CommentStorage interface {
	CreateComment(ctx context.Context, content string, postID int64, parentID *int64) error
	GetCommentByID(ctx context.Context, ID int64) (*models.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID int64) ([]models.Comment, error)
	DeleteComment(ctx context.Context, ID int64) error
}

type PostStorage interface {
	CreatePost(ctx context.Context, title, content string) error
	GetPostByID(ctx context.Context, ID int64) (*models.Post, error)
	DeletePost(ctx context.Context, ID int64) error
}
