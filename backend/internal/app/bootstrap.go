package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	"supermarket/internal/config"
	cancellationHandler "supermarket/internal/http/handlers/cancellation"
	priceHandler "supermarket/internal/http/handlers/price"
	productHandler "supermarket/internal/http/handlers/product"
	productSaleHandler "supermarket/internal/http/handlers/product_sale"
	productSupplyHandler "supermarket/internal/http/handlers/product_supply"
	producttypeHandler "supermarket/internal/http/handlers/product_type"
	saleHandler "supermarket/internal/http/handlers/sale"
	stockHandler "supermarket/internal/http/handlers/stock"
	supplierHandler "supermarket/internal/http/handlers/supplier"
	"supermarket/internal/http/router"
	"supermarket/internal/logger"
	cancellationRepo "supermarket/internal/repository/cancellation/gorm"
	priceRepo "supermarket/internal/repository/price/gorm"
	productRepo "supermarket/internal/repository/product/gorm"
	productSaleRepo "supermarket/internal/repository/product_sale/gorm"
	productSupplyRepo "supermarket/internal/repository/product_supply/gorm"
	producttypeRepo "supermarket/internal/repository/product_type/gorm"
	saleRepo "supermarket/internal/repository/sale/gorm"
	stockRepo "supermarket/internal/repository/stock/gorm"
	supplierRepo "supermarket/internal/repository/supplier/gorm"
	transactions "supermarket/internal/repository/transactions/gorm"
	cancellationService "supermarket/internal/services/cancellation"
	priceService "supermarket/internal/services/price"
	productService "supermarket/internal/services/product"
	productSaleService "supermarket/internal/services/product_sale"
	productSupplyService "supermarket/internal/services/product_supply"
	producttypeService "supermarket/internal/services/product_type"
	saleService "supermarket/internal/services/sale"
	stockService "supermarket/internal/services/stock"
	supplierService "supermarket/internal/services/supplier"
	"supermarket/internal/storage/postgres"
)

// New creates instant app.
func New(ctx context.Context) (*App, error) {
	// config
	if err := config.Init(); err != nil {
		slog.Error("config init failed", slog.Any("error", err))

		return nil, err
	}
	slog.Info("configuration loaded", slog.String("file", viper.ConfigFileUsed()))

	// logger
	logLevel := viper.GetString(config.LogLevel)
	l, err := logger.Init(logLevel)
	if err != nil {
		slog.Error("logger init failed", slog.Any("error", err))

		return nil, err
	}
	l.Info("logger initialized", slog.String("logLevel", logLevel))

	// storage
	dsn := config.GetPostgresDSN()
	db, err := postgres.New(ctx, l, dsn)
	if err != nil {
		l.Error("failed to connect db", slog.Any("error", err))

		return nil, err
	}
	l.Info("db connected")

	app := initApp(l, db)

	return &App{
		app: app,
		l:   l,
	}, nil
}

// initApp creates dependences, fiber app and setups router.
func initApp(l *slog.Logger, db *postgres.Storage) *fiber.App {
	// repositories
	unitOfWork := transactions.New(db.DB)
	productTypeR := producttypeRepo.New(db.DB)
	productR := productRepo.New(db.DB)
	priceR := priceRepo.New(db.DB)
	supplierR := supplierRepo.New(db.DB)
	stockR := stockRepo.New(db.DB)
	productSupplyR := productSupplyRepo.New(db.DB)
	cancellationR := cancellationRepo.New(db.DB)
	productSaleR := productSaleRepo.New(db.DB)
	saleR := saleRepo.New(db.DB)

	// services
	productTypeS := producttypeService.New(l, productTypeR)
	productS := productService.New(l, productR)
	priceS := priceService.New(l, priceR)
	supplierS := supplierService.New(l, supplierR)
	stockS := stockService.New(l, stockR)
	productSupplyS := productSupplyService.New(l, productSupplyR, unitOfWork)
	cancellationS := cancellationService.New(l, cancellationR, unitOfWork)
	productSaleS := productSaleService.New(l, productSaleR, unitOfWork)
	saleS := saleService.New(l, saleR, unitOfWork)

	readTimeout := time.Second * time.Duration(viper.GetInt(config.ReadTimeout))
	writeTimeout := time.Second * time.Duration(viper.GetInt(config.WriteTimeout))

	// handlers
	productTypeH := producttypeHandler.New(l, readTimeout, writeTimeout, productTypeS)
	productH := productHandler.New(l, readTimeout, writeTimeout, productS)
	priceH := priceHandler.New(l, readTimeout, writeTimeout, priceS)
	supplierH := supplierHandler.New(l, readTimeout, writeTimeout, supplierS)
	stockH := stockHandler.New(l, readTimeout, writeTimeout, stockS)
	productSupplyH := productSupplyHandler.New(l, readTimeout, writeTimeout, productSupplyS)
	cancellationH := cancellationHandler.New(l, readTimeout, writeTimeout, cancellationS)
	productSaleH := productSaleHandler.New(l, readTimeout, writeTimeout, productSaleS)
	saleH := saleHandler.New(l, readTimeout, writeTimeout, saleS)

	// app
	app := fiber.New(fiber.Config{})
	router.New(
		app,
		productTypeH,
		productH,
		priceH,
		supplierH,
		stockH,
		productSupplyH,
		cancellationH,
		productSaleH,
		saleH,
	).Setup()

	return app
}
