package main

import (
	"errors"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/viper"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"supermarket/internal/config"
)

func main() {
	if err := config.Init(); err != nil {
		slog.Error("config init failed", slog.Any("error", err))

		return
	}
	slog.Info("configuration loaded", slog.String("file", viper.ConfigFileUsed()))

	m, err := migrate.New(
		"file://migrations",
		config.GetPostgresDSN(),
	)
	if err != nil {
		slog.Error("failed to initialize migrations", slog.Any("error", err))

		return
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("no migrations to apply")

			return
		}

		slog.Error("failed to apply migrations", slog.Any("error", err))

		return
	}

	slog.Info("migrations applied successfully")
}
