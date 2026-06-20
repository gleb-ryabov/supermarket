package stock

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

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

// SetCountStock updates the quantity of a stock item by its product id. Sets += for count.
// If not found creates a new item.
func (s *service) SetCountStock(ctx context.Context, productID uuid.UUID, count decimal.Decimal) error {
	const op = "services.stocks.setCountStock"

	log := s.logger.
		With("op", op).
		With("productID", productID).
		With("count", count)

	if _, err := s.stocksR.FirstOrCreateByProductID(ctx, productID); err != nil {
		log.Error("failed to get or create stock",
			slog.Any("error", err),
		)

		return err
	}

	if err := s.stocksR.SetCountByProductID(ctx, productID, count); err != nil {
		if errors.Is(err, repository.ErrNotEnoughStock) {
			log.Error("Not enough stock quantity")

			return services.ErrNotEnoughStock
		}

		log.Error("failed to set count stock",
			slog.Any("error", err),
		)

		return err
	}

	return nil
}

// DeleteStock deletes stock in the db by product id.
func (s *service) DeleteStock(ctx context.Context, productID uuid.UUID) error {
	const op = "services.stocks.deleteStock"

	log := s.logger.With("op", op).
		With("productID", productID)

	stocks, err := s.GetStocks(ctx, "", &productID)
	if err != nil {
		return err
	}

	if len(stocks) == 0 {
		log.Warn("stock is not exists in the db", slog.Any("err", services.ErrNotFound))

		return nil
	}

	stock := stocks[0]

	if err = s.stocksR.Delete(ctx, stock.ID); err != nil {
		log.Error("failed to delete stock",
			slog.Any("error", err),
			slog.Any("id", stock.ID),
		)

		return err
	}

	return nil
}
