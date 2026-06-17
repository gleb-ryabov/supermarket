package supplier

import (
	"context"

	"github.com/google/uuid"

	"supermarket/internal/models"
)

// Service provides business logic for suppliers.
type Service interface {
	// GetSuppliers returns slice suppliers by search param.
	GetSuppliers(ctx context.Context, search string) ([]models.Supplier, error)
	// CreateSupplier creates supplier in the db.
	CreateSupplier(ctx context.Context, supplier *models.Supplier) error
	// DeleteSupplier deletes supplier in the db by id.
	DeleteSupplier(ctx context.Context, id uuid.UUID) error
	// UpdateSupplier updates supplier in the db.
	UpdateSupplier(ctx context.Context, supplier *models.Supplier) error
}
