package productsale

import (
	"context"

	"github.com/google/uuid"

	"supermarket/internal/models"
	"supermarket/internal/repository"
)

// Repository is the interface for product in the sale.
//
//go:generate mockery
type Repository interface {
	repository.Repository[models.ProductSale]
	// GetBySale returns product from sale.
	GetBySale(ctx context.Context, saleID uuid.UUID) ([]models.ProductSale, error)
	// DeleteBySale drops product by sale.
	DeleteBySale(ctx context.Context, saleID uuid.UUID) error
}
