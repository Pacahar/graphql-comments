package graphql

import (
	"log/slog"

	"github.com/Pacahar/graphql-comments/internal/graphql/generated"
	"github.com/Pacahar/graphql-comments/internal/storage"
)

type Resolver struct {
	Storage *storage.Storage
	Logger  *slog.Logger
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}
