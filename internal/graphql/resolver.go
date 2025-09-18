package graphql

import (
	"log/slog"

	"github.com/Pacahar/graphql-comments/internal/storage"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Storage *storage.Storage
	Logger  *slog.Logger
}
