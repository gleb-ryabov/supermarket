package stock

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"supermarket/internal/http/dto"
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
func (s *service) GetStocks(ctx context.Context, search string, productID *uuid.UUID) ([]dto.StockDTO, error) {
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

	result := make([]dto.StockDTO, 0, len(stocks))
	for _, v := range stocks {
		result = append(result, dto.ToStockDTO(&v))
	}

	return result, err
}

// IncreaseStock updates the quantity of a stock item by its product id. Sets += for count.
// If not found creates a new item.
func (s *service) IncreaseStock(
	ctx context.Context,
	productID uuid.UUID,
	count decimal.Decimal,
) error {
	return IncreaseStock(
		ctx,
		s.logger,
		s.stocksR,
		productID,
		count,
	)
}

// IncreaseStock updates the quantity of a stock item by its product id. Sets += for count.
// If not found creates a new item.
func IncreaseStock(
	ctx context.Context,
	logger *slog.Logger,
	stockR stock.Repository,
	productID uuid.UUID,
	count decimal.Decimal,
) error {
	const op = "services.stocks.setCountStock"

	log := logger.
		With("op", op).
		With("productID", productID).
		With("count", count)

	if _, err := stockR.FirstOrCreateByProductID(ctx, productID); err != nil {
		log.Error("failed to get or create stock",
			slog.Any("error", err),
		)

		return err
	}

	if err := stockR.SetCountByProductID(ctx, productID, count); err != nil {
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
