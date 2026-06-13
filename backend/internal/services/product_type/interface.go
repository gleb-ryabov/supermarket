package producttype

import (
	"context"
	"supermarket/internal/models"

	"github.com/google/uuid"
)

// Service  provides business logic for product types.
type Service interface {
	// GetProductTypes returns slice product types and error by params product name and for adult.
	GetProductTypes(ctx context.Context, name string, forAdult *bool) ([]models.ProductType, error)
	// CreateProductType creates product type in the db.
	CreateProductType(ctx context.Context, pt *models.ProductType) error
	// DeleteProductType deletes product type in the db by id.
	DeleteProductType(ctx context.Context, id uuid.UUID) error
	// CreateProductType updates product type in the db.
	UpdateProductType(ctx context.Context, pt *models.ProductType) error
}
