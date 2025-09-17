package memory

import (
	"context"
	"sync"
	"time"

	"github.com/Pacahar/graphql-comments/internal/models"
	storageErrors "github.com/Pacahar/graphql-comments/internal/storage/errors"
)

type PostMemoryStorage struct {
	mu        sync.RWMutex
	posts     map[int64]models.Post
	currentID int64
}

func (ps *PostMemoryStorage) CreatePost(ctx context.Context, title, content string, commentsDisabled bool) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.posts[ps.currentID] = models.Post{
		ID:               ps.currentID,
		Title:            title,
		Content:          content,
		CommentsDisabled: commentsDisabled,
		CreatedAt:        time.Now(),
	}

	ps.currentID++

	return nil
}

func (ps *PostMemoryStorage) GetPostByID(ctx context.Context, id int64) (models.Post, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	post, exists := ps.posts[id]
	if !exists {
		return models.Post{}, storageErrors.ErrPostNotFound
	}
	return post, nil
}

func (ps *PostMemoryStorage) DeletePost(ctx context.Context, id int64) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exists := ps.posts[id]; !exists {
		return storageErrors.ErrPostNotFound
	}

	delete(ps.posts, id)

	return nil
}

func NewPostMemoryStorage() (*PostMemoryStorage, error) {
	return &PostMemoryStorage{
		mu:        sync.RWMutex{},
		posts:     make(map[int64]models.Post),
		currentID: 1,
	}, nil
}
