package app

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

// App is struct for instant application.
type App struct {
	app *fiber.App
	l   *slog.Logger
}

// Run starts server.
func (a *App) Run(addr string) error {
	return a.app.Listen(addr)
}

// Shutdown stops server.
func (a *App) Shutdown(ctx context.Context) error {
	return a.app.ShutdownWithContext(ctx)
}
