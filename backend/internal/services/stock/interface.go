package stock

import (
	"context"

	"github.com/google/uuid"

	"supermarket/internal/models"
)

// Service  provides business logic for stocks.
type Service interface {
	// GetStocks returns slice stocks and error by params product id and for adult.
	GetStocks(ctx context.Context, search string, productID *uuid.UUID) ([]models.Stock, error)
}
