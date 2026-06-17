package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	"supermarket/internal/config"
	priceHandler "supermarket/internal/http/handlers/price"
	productHandler "supermarket/internal/http/handlers/product"
	producttypeHandler "supermarket/internal/http/handlers/product_type"
	supplierHandler "supermarket/internal/http/handlers/supplier"
	"supermarket/internal/http/router"
	"supermarket/internal/logger"
	priceRepo "supermarket/internal/repository/price/gorm"
	productRepo "supermarket/internal/repository/product/gorm"
	producttypeRepo "supermarket/internal/repository/product_type/gorm"
	supplierRepo "supermarket/internal/repository/supplier/gorm"
	priceService "supermarket/internal/services/price"
	productService "supermarket/internal/services/product"
	producttypeService "supermarket/internal/services/product_type"
	supplierService "supermarket/internal/services/supplier"
	"supermarket/internal/storage/postgres"
)

// New creates instant app.
func New(ctx context.Context) (*App, error) {
	//config
	if err := config.Init(); err != nil {
		slog.Error("config init failed", slog.Any("error", err))

		return nil, err
	}
	slog.Info("configuration loaded", slog.String("file", viper.ConfigFileUsed()))

	//logger
	logLevel := viper.GetString(config.LogLevel)
	l, err := logger.Init(logLevel)
	if err != nil {
		slog.Error("logger init failed", slog.Any("error", err))

		return nil, err
	}
	l.Info("logger initialized", slog.String("logLevel", logLevel))

	//storage
	dsn := config.GetPostgresDSN()
	db, err := postgres.New(ctx, l, dsn)
	if err != nil {
		l.Error("failed to connect db", slog.Any("error", err))

		return nil, err
	}
	l.Info("db connected")

	//repo
	productTypeR := producttypeRepo.New(db.DB)
	productR := productRepo.New(db.DB)
	priceR := priceRepo.New(db.DB)
	supplierR := supplierRepo.New(db.DB)

	//service
	productTypeS := producttypeService.New(l, productTypeR)
	productS := productService.New(l, productR)
	priceS := priceService.New(l, priceR)
	supplierS := supplierService.New(l, supplierR)

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
	supplierH := supplierHandler.New(
		l,
		time.Second*time.Duration(viper.GetInt(config.ReadTimeout)),
		time.Second*time.Duration(viper.GetInt(config.WriteTimeout)),
		supplierS,
	)

	//server
	app := fiber.New(fiber.Config{})
	router.New(app, productTypeH, productH, priceH, supplierH).Setup()

	return &App{
		app: app,
		l:   l,
	}, nil
}
