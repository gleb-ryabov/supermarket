package productsupply

import (
	"context"
	"time"

	"supermarket/internal/models"
	"supermarket/internal/repository"

	"github.com/google/uuid"
)

// Repository is the interface for product supplies.
type Repository interface {
	repository.Repository[models.ProductSupply]
	// GetByParams returns product supplies by params.
	GetByParams(
		ctx context.Context,
		productID *uuid.UUID,
		supplierId *uuid.UUID,
		dateFrom *time.Time,
		dateTo *time.Time,
	) ([]models.ProductSupply, error)
}
