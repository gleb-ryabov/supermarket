package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "supermarket/docs"
	"supermarket/internal/app"
	"supermarket/internal/config"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	appl, err := app.New(ctx)
	if err != nil {
		slog.Error("bootstrap failed", slog.Any("error", err))

		return
	}

	go func() {
		if err = appl.Run(config.GetServerURL()); err != nil {
			slog.Error("run server failed", slog.Any("error", err))
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received")
	if err = appl.Shutdown(ctx); err != nil {
		slog.Error("shutdown server failed", slog.Any("error", err))
	}
}
