package stock

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"supermarket/internal/models"
	"supermarket/internal/repository"
	"supermarket/internal/repository/stock"
	"supermarket/internal/services"
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

// UpdateCountStock updates the quantity of a stock item by its ID. Sets += for count.
func (s *service) UpdateCountStock(ctx context.Context, id uuid.UUID, count int) error {
	const op = "services.stocks.updateCountStock"

	log := s.logger.With("op", op)

	err := s.stocksR.UpdateCount(ctx, id, count)
	if err != nil {
		if errors.Is(err, repository.ErrNotEnoughStock) {
			log.Error("Not enough stock quantity",
				slog.Any("id", id),
				slog.Int("count", count),
			)

			return services.ErrNotEnoughStock
		}

		log.Error("failed to update count stock",
			slog.Any("error", err),
			slog.Any("id", id),
			slog.Int("count", count),
		)

		return err
	}

	return nil
}
