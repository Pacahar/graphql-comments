package memory

import "fmt"

type MemoryStorage struct {
	postStorage    *PostMemoryStorage
	commentStorage *CommentMemoryStorage
}

func NewMemoryStorage() (*MemoryStorage, error) {
	const op = "storage.memory.NewMemoryStorage"

	postStorage, err := NewPostMemoryStorage()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	commentStorage, err := NewCommentMemoryStorage()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &MemoryStorage{
		postStorage:    postStorage,
		commentStorage: commentStorage,
	}, nil
}
