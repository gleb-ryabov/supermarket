package stock

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"supermarket/internal/models"
	"supermarket/internal/repository/stock"
)

// service provides business logic for stocks.
type service struct {
	logger *slog.Logger

	stocksR stock.Repository
}

// New creates service for product types.
func New(
	logger *slog.Logger,
	stocksR stock.Repository,
) Service {
	return &service{
		logger:  logger,
		stocksR: stocksR,
	}
}

// GetStocks returns slice stocks and error by params product id and for adult.
func (s *service) GetStocks(ctx context.Context, search string, productID *uuid.UUID) ([]models.Stock, error) {
	const op = "services.stocks.getStocks"

	log := s.logger.With("op", op)

	stocks, err := s.stocksR.GetByParams(ctx, search, productID)
	if err != nil {
		log.Error("failed to get product types",
			slog.Any("error", err),
			slog.String("search", search),
			slog.Any("productId", productID),
		)

		return nil, err
	}

	return stocks, err
}
