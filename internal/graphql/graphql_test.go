package graphql

import (
	"context"
	"strconv"
	"testing"

	"log/slog"

	"github.com/Pacahar/graphql-comments/internal/storage"
	"github.com/Pacahar/graphql-comments/internal/storage/memory"
	"github.com/stretchr/testify/assert"
)

func setupResolver(t *testing.T) *Resolver {
	postStorage, err := memory.NewPostMemoryStorage()
	assert.NoError(t, err)

	commentStorage, err := memory.NewCommentMemoryStorage()
	assert.NoError(t, err)

	resolver := &Resolver{
		Storage: &storage.Storage{
			Post:    postStorage,
			Comment: commentStorage,
		},
		Logger: slog.New(slog.NewTextHandler(&testWriter{}, &slog.HandlerOptions{})),
	}

	return resolver
}

func TestCreateAndFetchPost(t *testing.T) {
	resolver := setupResolver(t)
	ctx := context.Background()
	mutation := &mutationResolver{resolver}

	post, err := mutation.CreatePost(ctx, "Title", "Content", false)
	assert.NoError(t, err)
	assert.Equal(t, "Title", post.Title)
	assert.Equal(t, "Content", post.Content)
	assert.False(t, post.CommentsDisabled)

	query := &queryResolver{resolver}
	fetched, err := query.Post(ctx, post.ID)
	assert.NoError(t, err)
	assert.Equal(t, post.ID, fetched.ID)
	assert.Len(t, fetched.Comments, 0)
}

func TestCreateComment(t *testing.T) {
	resolver := setupResolver(t)
	ctx := context.Background()
	mutation := &mutationResolver{resolver}

	post, _ := mutation.CreatePost(ctx, "Title", "Content", false)

	comment, err := mutation.CreateComment(ctx, post.ID, "comment", nil)
	assert.NoError(t, err)
	assert.Equal(t, "comment", comment.Content)
	assert.Equal(t, post.ID, comment.PostID)
	assert.Nil(t, comment.ParentID)

	childComment, err := mutation.CreateComment(ctx, post.ID, "Child comment", &comment.ID)
	assert.NoError(t, err)
	assert.Equal(t, comment.ID, *childComment.ParentID)

	query := &queryResolver{resolver}
	fetchedComments, err := query.Comments(ctx, post.ID, nil, nil)
	assert.NoError(t, err)
	assert.Len(t, fetchedComments, 2)

	var foundChild bool
	for _, c := range fetchedComments {
		if c.ParentID != nil && *c.ParentID == comment.ID {
			foundChild = true
		}
	}
	assert.True(t, foundChild)
}

func TestDeletePostAndComments(t *testing.T) {
	resolver := setupResolver(t)
	ctx := context.Background()
	mutation := &mutationResolver{resolver}

	post, _ := mutation.CreatePost(ctx, "Post", "Content", false)
	comment, _ := mutation.CreateComment(ctx, post.ID, "Comment", nil)

	ok, err := mutation.DeletePost(ctx, post.ID)
	assert.NoError(t, err)
	assert.True(t, ok)

	query := &queryResolver{resolver}
	_, err = query.Post(ctx, post.ID)
	assert.Error(t, err)

	commentID, _ := strconv.ParseInt(comment.ID, 10, 64)
	_, err = resolver.Storage.Comment.GetCommentByID(ctx, commentID)
	assert.Error(t, err)
}

func TestDeleteComment(t *testing.T) {
	resolver := setupResolver(t)
	ctx := context.Background()
	mutation := &mutationResolver{resolver}

	post, _ := mutation.CreatePost(ctx, "Post", "Content", false)
	comment, _ := mutation.CreateComment(ctx, post.ID, "Comment", nil)

	ok, err := mutation.DeleteComment(ctx, comment.ID)
	assert.NoError(t, err)
	assert.True(t, ok)

	commentID, _ := strconv.ParseInt(comment.ID, 10, 64)
	_, err = resolver.Storage.Comment.GetCommentByID(ctx, commentID)
	assert.Error(t, err)
}

func TestFetchPostsWithPagination(t *testing.T) {
	resolver := setupResolver(t)
	ctx := context.Background()
	mutation := &mutationResolver{resolver}
	query := &queryResolver{resolver}

	for i := 1; i <= 5; i++ {
		_, _ = mutation.CreatePost(ctx, "Post "+strconv.Itoa(i), "Content", false)
	}

	posts, err := query.Posts(ctx, nil, nil)
	assert.NoError(t, err)
	assert.Len(t, posts, 5)

	limit := int32(2)
	offset := int32(1)
	paged, err := query.Posts(ctx, &limit, &offset)
	assert.NoError(t, err)
	assert.Len(t, paged, 2)
}

type testWriter struct{}

func (tw *testWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
