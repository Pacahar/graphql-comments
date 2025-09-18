package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Pacahar/graphql-comments/internal/config"
	"github.com/Pacahar/graphql-comments/internal/constants"
	"github.com/Pacahar/graphql-comments/internal/graphql"
	"github.com/Pacahar/graphql-comments/internal/graphql/generated"
	"github.com/Pacahar/graphql-comments/internal/storage"
	"github.com/Pacahar/graphql-comments/internal/storage/memory"
	"github.com/Pacahar/graphql-comments/internal/storage/postgres"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Environment)

	log.Info("Starting service", slog.String("env", cfg.Environment))
	log.Debug("Debug messages enabled")

	storage, err := setupStorage(&cfg.Storage)

	if err != nil {
		log.Error("failed to setup storage", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("storage set", slog.String("storage type", cfg.Storage.Type))

	resolver := &graphql.Resolver{
		Storage: storage,
		Logger:  log,
	}

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: resolver}),
	)

	http.Handle("/playground", playground.Handler("GraphQL playground", "/query"))

	http.Handle("/query", srv)

	address := fmt.Sprintf(":%d", cfg.HTTPServer.Port)
	log.Info("Starting GraphQL server", slog.Int("addr", cfg.HTTPServer.Port))

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Error("failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case constants.EnvLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case constants.EnvDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case constants.EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupStorage(storageCfg *config.Storage) (*storage.Storage, error) {
	switch storageCfg.Type {
	case constants.StorageMemory:
		return memory.NewMemoryStorage()
	case constants.StoragePostgres:
		return postgres.NewPostgresStorage(storageCfg.Postgres.DSN())
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageCfg.Type)
	}
}
