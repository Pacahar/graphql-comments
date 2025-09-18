package graphql

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/Pacahar/graphql-comments/internal/graphql/generated"
)

type mutationResolver struct{ *Resolver }

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, commentsDisabled bool) (*generated.Post, error) {
	id, err := r.Storage.Post.CreatePost(ctx, title, content, commentsDisabled)

	if err != nil {
		r.Logger.Error("failed to create post", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to create post")
	}

	r.Logger.Info("post created successfully", slog.Int64("id", id))

	post, err := r.Storage.Post.GetPostByID(ctx, id)

	if err != nil {
		r.Logger.Error("failed to fetch created post", slog.String("err", err.Error()))
		return nil, fmt.Errorf("internal error")
	}

	// comments, err := r.Storage.Comment.GetCommentsByPostID(ctx, id)

	gqlComments := make([]*generated.Comment, 0)

	// Here comments always gonna be empty so we don't need this

	// for _, comment := range comments {

	// 	var parentIDCopy *string
	// 	if comment.ParentID != nil {
	// 		s := strconv.FormatInt(*comment.ParentID, 10)
	// 		parentIDCopy = &s
	// 	}

	// 	gqlComments = append(gqlComments, &Comment{
	// 		ID:        strconv.FormatInt(comment.ID, 10),
	// 		PostID:    strconv.FormatInt(comment.PostID, 10),
	// 		ParentID:  parentIDCopy,
	// 		Content:   comment.Content,
	// 		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
	// 	})
	// }

	return &generated.Post{
		ID:               strconv.FormatInt(post.ID, 10),
		Title:            post.Title,
		Content:          post.Content,
		CommentsDisabled: post.CommentsDisabled,
		CreatedAt:        post.CreatedAt.Format(time.RFC3339),
		Comments:         gqlComments,
	}, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, postID string, content string, parentID *string) (*generated.Comment, error) {
	var pInt64ParentID *int64

	if parentID != nil {
		intParentID, err := strconv.Atoi(*parentID)

		if err != nil {
			r.Logger.Error("invalid parent id", slog.String("err", err.Error()), slog.String("id", *parentID))
			return nil, fmt.Errorf("invalid parent id ")
		}

		int64ParentID := int64(intParentID)

		_, err = r.Storage.Comment.GetCommentByID(ctx, int64ParentID)

		if err != nil {
			r.Logger.Error("failed to fetch parent comment", slog.String("err", err.Error()))
			return nil, fmt.Errorf("failed to fetch parent comment")
		}

		pInt64ParentID = &int64ParentID
	}

	intPostID, err := strconv.Atoi(postID)

	if err != nil {
		r.Logger.Error("invalid post id", slog.String("err", err.Error()), slog.String("id", postID))
		return nil, fmt.Errorf("invalid post id")
	}

	post, err := r.Storage.Post.GetPostByID(ctx, int64(intPostID))

	if err != nil {
		r.Logger.Error("failed to fetch post", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to fetch post")
	}

	if post.CommentsDisabled {
		r.Logger.Error("comments disabled on this post", slog.String("err", err.Error()))
		return nil, fmt.Errorf("comments disabled on this post")
	}

	id, err := r.Storage.Comment.CreateComment(ctx, content, int64(intPostID), pInt64ParentID)

	if err != nil {
		r.Logger.Error("failed to create comment")
		return nil, fmt.Errorf("failed to create comment")
	}

	r.Logger.Info("comment created successfully", slog.Int64("id", id))

	comment, err := r.Storage.Comment.GetCommentByID(ctx, id)

	if err != nil {
		r.Logger.Error("failed to fetch created comment", slog.String("err", err.Error()))
		return nil, fmt.Errorf("internal error")
	}

	var parentIDCopy *string
	if comment.ParentID != nil {
		s := strconv.FormatInt(*comment.ParentID, 10)
		parentIDCopy = &s
	}

	return &generated.Comment{
		ID:        strconv.FormatInt(comment.ID, 10),
		PostID:    strconv.FormatInt(comment.PostID, 10),
		ParentID:  parentIDCopy,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
	}, nil
}

// DeletePost is the resolver for the deletePost field.
func (r *mutationResolver) DeletePost(ctx context.Context, id string) (bool, error) {
	intID, err := strconv.Atoi(id)

	if err != nil {
		r.Logger.Error("invalid post id", slog.String("err", err.Error()))
		return false, fmt.Errorf("invalid post id")
	}

	_, err = r.Storage.Post.GetPostByID(ctx, int64(intID))

	if err != nil {
		r.Logger.Error("post not found", slog.String("err", err.Error()), slog.Int("id", intID))
		return false, fmt.Errorf("post not found")
	}

	err = r.Storage.Comment.DeleteCommentsByPostID(ctx, int64(intID))

	if err != nil {
		r.Logger.Error("failed to delete comments from post", slog.String("err", err.Error()))
		return false, fmt.Errorf("failed to delete comments from post")
	}

	err = r.Storage.Post.DeletePost(ctx, int64(intID))

	if err != nil {
		r.Logger.Error("failed to delete post", slog.String("err", err.Error()))
		return false, fmt.Errorf("failed to delete post")
	}

	r.Logger.Info("post deleted successfully")

	return true, nil
}

// DeleteComment is the resolver for the deleteComment field.
func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (bool, error) {
	intID, err := strconv.Atoi(id)

	if err != nil {
		r.Logger.Error("invalid comment ID", slog.String("err", err.Error()))
		return false, fmt.Errorf("invalid comment ID")
	}

	_, err = r.Storage.Comment.GetCommentByID(ctx, int64(intID))

	if err != nil {
		r.Logger.Error("comment not found", slog.String("err", err.Error()))
		return false, fmt.Errorf("comment not found")
	}

	err = r.Storage.Comment.DeleteComment(ctx, int64(intID))

	if err != nil {
		r.Logger.Error("failed to delete comment", slog.String("err", err.Error()))
		return false, fmt.Errorf("failed to delete comment")
	}

	r.Logger.Info("comment deleted successfully")

	return true, nil
}
