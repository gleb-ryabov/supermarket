package productsupply

import (
	"context"
	"time"

	"github.com/google/uuid"

	"supermarket/internal/models"
	"supermarket/internal/repository"
)

// Repository is the interface for product supplies.
type Repository interface {
	repository.Repository[models.ProductSupply]
	// GetByParams returns product supplies by params.
	GetByParams(
		ctx context.Context,
		productID *uuid.UUID,
		supplierID *uuid.UUID,
		dateFrom *time.Time,
		dateTo *time.Time,
	) ([]models.ProductSupply, error)
}
