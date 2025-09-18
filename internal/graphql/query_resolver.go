package graphql

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"
)

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*Post, error) {
	intID, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		r.Logger.Error("invalid post id", slog.String("err", err.Error()))
		return nil, fmt.Errorf("invalid post id")
	}

	post, err := r.Storage.Post.GetPostByID(ctx, intID)

	if err != nil {
		r.Logger.Error("failed to fetch post", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to fetch post")
	}

	comments, err := r.Storage.Comment.GetCommentsByPostID(ctx, intID, nil, nil)

	if err != nil {
		r.Logger.Error("failed to fetch comments", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to fetch comments")
	}

	gqlComments := make([]*Comment, 0, len(comments))

	for _, comment := range comments {
		childComments, err := r.Storage.Comment.GetCommentsByParentID(ctx, comment.ID)

		if err != nil {
			r.Logger.Error("failed to fetch child comments", slog.String("err", err.Error()))
			return nil, fmt.Errorf("failed to fetch child comments")
		}

		gqlReplies := make([]*Comment, 0, len(childComments))

		for _, child := range childComments {
			var childParentIDCopy *string
			if child.ParentID != nil {
				s := strconv.FormatInt(*child.ParentID, 10)
				childParentIDCopy = &s
			}

			gqlReplies = append(gqlReplies, &Comment{
				strconv.FormatInt(child.ID, 10),
				strconv.FormatInt(child.PostID, 10),
				childParentIDCopy,
				child.Content,
				child.CreatedAt.Format(time.RFC3339),
				nil,
			})
		}

		var parentIDCopy *string
		if comment.ParentID != nil {
			s := strconv.FormatInt(*comment.ParentID, 10)
			parentIDCopy = &s
		}

		gqlComments = append(gqlComments, &Comment{
			strconv.FormatInt(comment.ID, 10),
			strconv.FormatInt(comment.PostID, 10),
			parentIDCopy,
			comment.Content,
			comment.CreatedAt.Format(time.RFC3339),
			gqlReplies,
		})
	}

	return &Post{
		strconv.FormatInt(post.ID, 10),
		post.Title,
		post.Content,
		post.CommentsDisabled,
		post.CreatedAt.Format(time.RFC3339),
		gqlComments,
	}, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, limit *int32, offset *int32) ([]*Post, error) {
	posts, err := r.Storage.Post.GetAllPosts(ctx)

	if err != nil {
		r.Logger.Error("failed to fetch posts", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to fetch posts")
	}

	gqlPosts := make([]*Post, 0, len(posts))

	for _, post := range posts {
		gqlPosts = append(gqlPosts, &Post{
			strconv.FormatInt(post.ID, 10),
			post.Title,
			post.Content,
			post.CommentsDisabled,
			post.CreatedAt.Format(time.RFC3339),
			nil,
		})
	}

	r.Logger.Info("Fetch all posts successfully")

	if limit != nil && offset != nil {
		start := int(*offset)
		end := start + int(*limit)
		if start > len(gqlPosts) {
			return []*Post{}, nil
		}
		if end > len(gqlPosts) {
			end = len(gqlPosts)
		}
		gqlPosts = gqlPosts[start:end]
	}

	return gqlPosts, nil
}

// Comment is the resolver for the comment field.
func (r *queryResolver) Comment(ctx context.Context, id string) (*Comment, error) {
	intID, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		r.Logger.Error("invalid comment id", slog.String("err", err.Error()))
		return nil, fmt.Errorf("invalid comment id")
	}

	comment, err := r.Storage.Comment.GetCommentByID(ctx, intID)

	if err != nil {
		r.Logger.Error("failed to fetch comment", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to fetch comment")
	}

	childComments, err := r.Storage.Comment.GetCommentsByParentID(ctx, comment.ID)

	if err != nil {
		r.Logger.Error("failed to fetch child comments", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to fetch child comments")
	}

	gqlReplies := make([]*Comment, 0, len(childComments))

	for _, child := range childComments {
		gqlReplies = append(gqlReplies, &Comment{
			strconv.FormatInt(child.ID, 10),
			strconv.FormatInt(child.PostID, 10),
			&id,
			child.Content,
			child.CreatedAt.Format(time.RFC3339),
			nil,
		})
	}

	var parentIDCopy *string
	if comment.ParentID != nil {
		s := strconv.FormatInt(*comment.ParentID, 10)
		parentIDCopy = &s
	}

	return &Comment{
		strconv.FormatInt(comment.ID, 10),
		strconv.FormatInt(comment.PostID, 10),
		parentIDCopy,
		comment.Content,
		comment.CreatedAt.Format(time.RFC3339),
		gqlReplies,
	}, nil
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID string, limit *int64, offset *int64) ([]*Comment, error) {
	intPostID, err := strconv.ParseInt(postID, 10, 64)

	if err != nil {
		r.Logger.Error("invalid post id", slog.String("err", err.Error()))
		return nil, fmt.Errorf("invalid post id")
	}

	comments, err := r.Storage.Comment.GetCommentsByPostID(ctx, intPostID, limit, offset)

	if err != nil {
		r.Logger.Error("failed to fetch comments", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to fetch comments")
	}

	gqlComments := make([]*Comment, 0, len(comments))

	for _, comment := range comments {
		childComments, err := r.Storage.Comment.GetCommentsByParentID(ctx, comment.ID)

		if err != nil {
			r.Logger.Error("failed to fetch child comments", slog.String("err", err.Error()))
			return nil, fmt.Errorf("failed to fetch child comments")
		}

		gqlChildComments := make([]*Comment, 0, len(childComments))

		for _, child := range childComments {
			var childParentIDCopy *string
			if child.ParentID != nil {
				s := strconv.FormatInt(*child.ParentID, 10)
				childParentIDCopy = &s
			}

			gqlChildComments = append(gqlChildComments, &Comment{
				strconv.FormatInt(child.ID, 10),
				strconv.FormatInt(child.PostID, 10),
				childParentIDCopy,
				child.Content,
				child.CreatedAt.Format(time.RFC3339),
				nil,
			})
		}

		var parentIDCopy *string
		if comment.ParentID != nil {
			s := strconv.FormatInt(*comment.ParentID, 10)
			parentIDCopy = &s
		}

		gqlComments = append(gqlComments, &Comment{
			strconv.FormatInt(comment.ID, 10),
			strconv.FormatInt(comment.PostID, 10),
			parentIDCopy,
			comment.Content,
			comment.CreatedAt.Format(time.RFC3339),
			gqlChildComments,
		})
	}

	return gqlComments, nil
}
