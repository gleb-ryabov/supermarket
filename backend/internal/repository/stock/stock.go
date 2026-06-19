package stock

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"supermarket/internal/models"
	"supermarket/internal/repository"
)

// Repository is the interface for stocks.
type Repository interface {
	repository.Repository[models.Stock]
	// GetByParams returns stocks by search param.
	GetByParams(ctx context.Context, search string, productID *uuid.UUID) ([]models.Stock, error)
	// SetCountByProductID updates the quantity of a stock item by product ID.
	SetCountByProductID(ctx context.Context, productID uuid.UUID, count decimal.Decimal) error
	// FirstOrCreateByProductID finds the stock by product ID, otherwise if not found creates a new stock.
	FirstOrCreateByProductID(ctx context.Context, productID uuid.UUID) (*models.Stock, error)
}
