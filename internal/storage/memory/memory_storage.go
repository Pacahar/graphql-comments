package memory

import (
	"fmt"

	"github.com/Pacahar/graphql-comments/internal/storage"
)

func NewMemoryStorage() (*storage.Storage, error) {
	const op = "storage.memory.NewMemoryStorage"

	postStorage, err := NewPostMemoryStorage()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	commentStorage, err := NewCommentMemoryStorage()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &storage.Storage{
		Post:    postStorage,
		Comment: commentStorage,
	}, nil
}
