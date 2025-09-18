package memory

import (
	"context"
	"sync"
	"time"

	"github.com/Pacahar/graphql-comments/internal/models"
	storageErrors "github.com/Pacahar/graphql-comments/internal/storage/errors"
)

type CommentMemoryStorage struct {
	mu        sync.RWMutex
	comments  map[int64]models.Comment
	currentID int64
}

func NewCommentMemoryStorage() (*CommentMemoryStorage, error) {
	return &CommentMemoryStorage{
		mu:        sync.RWMutex{},
		comments:  make(map[int64]models.Comment),
		currentID: 1,
	}, nil
}

func (cs *CommentMemoryStorage) CreateComment(ctx context.Context, content string, postID int64, parentID *int64) (int64, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	id := cs.currentID

	var safeParentID *int64
	if parentID != nil {
		val := *parentID
		safeParentID = &val
	}

	cs.comments[id] = models.Comment{
		ID:        id,
		PostID:    postID,
		ParentID:  safeParentID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	cs.currentID++

	return id, nil
}

func (cs *CommentMemoryStorage) GetCommentByID(ctx context.Context, id int64) (models.Comment, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	comment, exists := cs.comments[id]
	if !exists {
		return models.Comment{}, storageErrors.ErrCommentNotFound
	}

	var parentID *int64
	if comment.ParentID != nil {
		val := *comment.ParentID
		parentID = &val
	}

	copy := models.Comment{
		ID:        comment.ID,
		PostID:    comment.PostID,
		ParentID:  parentID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}

	return copy, nil
}

func (cs *CommentMemoryStorage) GetCommentsByPostID(ctx context.Context, postID int64) ([]models.Comment, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	resultComments := make([]models.Comment, 0)

	for _, comment := range cs.comments {
		if comment.PostID == postID {

			var parentID *int64
			if comment.ParentID != nil {
				val := *comment.ParentID
				parentID = &val
			}

			copy := models.Comment{
				ID:        comment.ID,
				PostID:    comment.PostID,
				ParentID:  parentID,
				Content:   comment.Content,
				CreatedAt: comment.CreatedAt,
			}

			resultComments = append(resultComments, copy)
		}
	}

	return resultComments, nil
}

func (cs *CommentMemoryStorage) DeleteComment(ctx context.Context, id int64) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	var deleteRecursive func(commentID int64)

	deleteRecursive = func(commentID int64) {
		for childID, comment := range cs.comments {

			if comment.ParentID != nil && *comment.ParentID == commentID {
				deleteRecursive(childID)
			}
		}
		delete(cs.comments, commentID)
	}

	deleteRecursive(id)

	return nil
}

func (cs *CommentMemoryStorage) DeleteCommentsByPostID(ctx context.Context, postID int64) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	for id, comment := range cs.comments {
		if comment.PostID == postID {
			delete(cs.comments, id)
		}
	}

	return nil
}
