package main

import (
	"log/slog"
	"os"

	"github.com/Pacahar/graphql-comments/internal/config"
)

const (
	envLocal string = "local"
	envDev   string = "dev"
	envProd  string = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Environment)

	log.Info("Starting service", slog.String("env", cfg.Environment))
	log.Debug("Debug messages enabled")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
