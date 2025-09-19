package memory

import (
	"context"
	"testing"

	storageErrors "github.com/Pacahar/graphql-comments/internal/storage/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetPost(t *testing.T) {
	ctx := context.Background()

	storage, err := NewPostMemoryStorage()
	assert.NoError(t, err)

	id, err := storage.CreatePost(ctx, "Title 1", "Content 1", false)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)

	post, err := storage.GetPostByID(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, "Title 1", post.Title)
	assert.Equal(t, "Content 1", post.Content)
	assert.False(t, post.CommentsDisabled)
}

func TestGetAllPosts(t *testing.T) {
	ctx := context.Background()

	storage, err := NewPostMemoryStorage()
	assert.NoError(t, err)

	_, err = storage.CreatePost(ctx, "Post1", "Content1", false)
	assert.NoError(t, err)

	_, err = storage.CreatePost(ctx, "Post2", "Content2", true)
	assert.NoError(t, err)

	posts, err := storage.GetAllPosts(ctx)
	assert.NoError(t, err)
	assert.Len(t, posts, 2)
}

func TestDeletePost(t *testing.T) {
	ctx := context.Background()
	storage, err := NewPostMemoryStorage()
	assert.NoError(t, err)

	id, err := storage.CreatePost(ctx, "Title", "Content", false)
	assert.NoError(t, err)

	err = storage.DeletePost(ctx, id)
	assert.NoError(t, err)

	_, err = storage.GetPostByID(ctx, id)
	assert.ErrorIs(t, err, storageErrors.ErrPostNotFound)
}

func TestCreateAndGetComment(t *testing.T) {
	ctx := context.Background()

	storage, err := NewCommentMemoryStorage()
	assert.NoError(t, err)

	postID := int64(1)
	commentID, err := storage.CreateComment(ctx, "Comment 1", postID, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), commentID)

	comment, err := storage.GetCommentByID(ctx, commentID)
	assert.NoError(t, err)
	assert.Equal(t, "Comment 1", comment.Content)
	assert.Equal(t, postID, comment.PostID)
	assert.Nil(t, comment.ParentID)
}

func TestGetCommentsByPostID(t *testing.T) {
	ctx := context.Background()

	storage, err := NewCommentMemoryStorage()
	assert.NoError(t, err)

	postID := int64(1)

	_, err = storage.CreateComment(ctx, "Comment 1", postID, nil)
	assert.NoError(t, err)

	_, err = storage.CreateComment(ctx, "Comment 2", postID, nil)
	assert.NoError(t, err)

	comments, err := storage.GetCommentsByPostID(ctx, postID, nil, nil)
	assert.NoError(t, err)
	assert.Len(t, comments, 2)
}

func TestGetCommentsByParentID(t *testing.T) {
	ctx := context.Background()

	storage, err := NewCommentMemoryStorage()
	assert.NoError(t, err)

	postID := int64(1)
	parentID, err := storage.CreateComment(ctx, "Parent", postID, nil)
	assert.NoError(t, err)

	childID, err := storage.CreateComment(ctx, "Child", postID, &parentID)
	assert.NoError(t, err)

	children, err := storage.GetCommentsByParentID(ctx, parentID)
	assert.NoError(t, err)

	assert.Len(t, children, 1)
	assert.Equal(t, childID, children[0].ID)
}

func TestDeleteComment(t *testing.T) {
	ctx := context.Background()
	storage, err := NewCommentMemoryStorage()
	assert.NoError(t, err)

	postID := int64(1)
	parentID, err := storage.CreateComment(ctx, "Parent", postID, nil)
	assert.NoError(t, err)

	_, err = storage.CreateComment(ctx, "Child", postID, &parentID)
	assert.NoError(t, err)

	err = storage.DeleteComment(ctx, parentID)
	assert.NoError(t, err)

	_, err = storage.GetCommentByID(ctx, parentID)
	assert.ErrorIs(t, err, storageErrors.ErrCommentNotFound)

	children, _ := storage.GetCommentsByParentID(ctx, parentID)
	assert.Len(t, children, 0)
}

func TestDeleteCommentsByPostID(t *testing.T) {
	ctx := context.Background()
	CommentStorage, _ := NewCommentMemoryStorage()
	PostStorage, _ := NewPostMemoryStorage()

	postID, err := PostStorage.CreatePost(ctx, "Post", "Content", false)
	assert.NoError(t, err)

	_, err = CommentStorage.CreateComment(ctx, "Comment 1", postID, nil)
	assert.NoError(t, err)

	_, err = CommentStorage.CreateComment(ctx, "Comment 2", postID, nil)
	assert.NoError(t, err)

	err = CommentStorage.DeleteCommentsByPostID(ctx, postID)
	assert.NoError(t, err)

	comments, _ := CommentStorage.GetCommentsByPostID(ctx, postID, nil, nil)
	assert.Len(t, comments, 0)
}
