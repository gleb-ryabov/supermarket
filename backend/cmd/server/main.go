package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	"supermarket/internal/config"
	priceHandler "supermarket/internal/http/handlers/price"
	productHandler "supermarket/internal/http/handlers/product"
	producttypeHandler "supermarket/internal/http/handlers/product_type"
	"supermarket/internal/http/router"
	"supermarket/internal/logger"
	priceRepo "supermarket/internal/repository/price/gorm"
	productRepo "supermarket/internal/repository/product/gorm"
	producttypeRepo "supermarket/internal/repository/product_type/gorm"
	priceService "supermarket/internal/services/price"
	productService "supermarket/internal/services/product"
	producttypeService "supermarket/internal/services/product_type"
	"supermarket/internal/storage/postgres"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	//config
	if err := config.Init(); err != nil {
		slog.Error("config init failed", slog.Any("error", err))

		return
	}
	slog.Info("configuration loaded", slog.String("file", viper.ConfigFileUsed()))

	//logger
	logLevel := viper.GetString(config.LogLevel)
	l, err := logger.Init(logLevel)
	if err != nil {
		slog.Error("logger init failed", slog.Any("error", err))

		return
	}
	l.Info("logger initialized", slog.String("logLevel", logLevel))

	//storage
	dsn := config.GetPostgresDSN()
	db, err := postgres.New(ctx, l, dsn)
	if err != nil {
		l.Error("failed to connect db", slog.Any("error", err))

		return
	}
	l.Info("db connected")

	//repo
	productTypeR := producttypeRepo.New(db.DB)
	productR := productRepo.New(db.DB)
	priceR := priceRepo.New(db.DB)

	//service
	productTypeS := producttypeService.New(l, productTypeR)
	productS := productService.New(l, productR)
	priceS := priceService.New(l, priceR)

	//handler
	productTypeH := producttypeHandler.New(
		l,
		time.Second*time.Duration(viper.GetInt(config.ReadTimeout)),
		time.Second*time.Duration(viper.GetInt(config.WriteTimeout)),
		productTypeS,
	)
	productH := productHandler.New(
		l,
		time.Second*time.Duration(viper.GetInt(config.ReadTimeout)),
		time.Second*time.Duration(viper.GetInt(config.WriteTimeout)),
		productS,
	)
	priceH := priceHandler.New(
		l,
		time.Second*time.Duration(viper.GetInt(config.ReadTimeout)),
		time.Second*time.Duration(viper.GetInt(config.WriteTimeout)),
		priceS,
	)

	//server
	app := fiber.New(fiber.Config{})
	router.New(app, productTypeH, productH, priceH).Setup()

	adr := config.GetServerURL()
	go func() {
		if err = app.Listen(adr); err != nil {
			l.Error("fiber server failed", slog.Any("error", err))

			return
		}
	}()
	l.Info("server started")

	<-ctx.Done()
	l.Info("shutdown signal received")
	_ = app.Shutdown()
}
