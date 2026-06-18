package stock

import (
	"context"

	"github.com/google/uuid"

	"supermarket/internal/models"
	"supermarket/internal/repository"
)

// Repository is the interface for stocks.
type Repository interface {
	repository.Repository[models.Stock]
	// GetByParams returns stocks by search param.
	GetByParams(ctx context.Context, search string, productID *uuid.UUID) ([]models.Stock, error)
	// UpdateCount updates the quantity of a stock item by its ID.
	UpdateCount(ctx context.Context, id uuid.UUID, count int) error
}
