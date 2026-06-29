package stock

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"supermarket/internal/http/dto"
)

// Service  provides business logic for stocks.
type Service interface {
	// GetStocks returns slice stocks and error by params product id and for adult.
	GetStocks(ctx context.Context, search string, productID *uuid.UUID) ([]dto.StockDTO, error)
	// IncreaseStock updates the quantity of a stock item by its product ID. Sets += for count.
	IncreaseStock(ctx context.Context, productID uuid.UUID, count decimal.Decimal) error
}
