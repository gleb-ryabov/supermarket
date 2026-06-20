package productsupply

import (
	"context"
	"time"

	"github.com/google/uuid"

	"supermarket/internal/http/dto"
	"supermarket/internal/models"
)

// Service provides business logic for product supplies.
type Service interface {
	// GetProductSupplies returns slice product supplies and error by params product id, supplie id and period delivery.
	GetProductSupplies(
		ctx context.Context,
		productID *uuid.UUID,
		supplierID *uuid.UUID,
		dateFrom *time.Time,
		dateTo *time.Time,
	) ([]dto.ProductSupplyDTO, error)
	// CreateProductSupply creates product supply in the db.
	CreateProductSupply(ctx context.Context, ps *models.ProductSupply) error
	// DeleteProductSupply deletes product supply in the db by id.
	DeleteProductSupply(ctx context.Context, id uuid.UUID) error
	// UpdateProductSupply deletes product supply in the db by id.
	UpdateProductSupply(ctx context.Context, ps *models.ProductSupply) error
}
