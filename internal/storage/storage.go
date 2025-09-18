package storage

import (
	"context"

	"github.com/Pacahar/graphql-comments/internal/models"
)

type Storage struct {
	Post    PostStorage
	Comment CommentStorage
}

type PostStorage interface {
	CreatePost(ctx context.Context, title, content string, commentsDisabled bool) (int64, error)
	GetPostByID(ctx context.Context, id int64) (models.Post, error)
	GetAllPosts(ctx context.Context) ([]models.Post, error)
	DeletePost(ctx context.Context, id int64) error
}

type CommentStorage interface {
	CreateComment(ctx context.Context, content string, postID int64, parentID *int64) (int64, error)
	GetCommentByID(ctx context.Context, id int64) (models.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID int64, limit *int64, offset *int64) ([]models.Comment, error)
	DeleteComment(ctx context.Context, id int64) error
	DeleteCommentsByPostID(ctx context.Context, id int64) error
}
