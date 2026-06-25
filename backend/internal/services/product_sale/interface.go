package productsale

import (
	"context"

	"github.com/google/uuid"

	"supermarket/internal/http/dto"
	"supermarket/internal/models"
)

// Service provides business logic for product in sale.
type Service interface {
	// GetProductsInSale returns slice products from sale.
	GetProductsInSale(ctx context.Context, saleID uuid.UUID) ([]dto.ProductSaleDTO, error)
	// CreateProductInSale adds a product to a sale.
	CreateProductInSale(ctx context.Context, productSale *models.ProductSale) error
	// UpdateProductInSale updates a product to a sale.
	UpdateProductInSale(ctx context.Context, productSale *models.ProductSale) error
	// DeleteProductInSale drops a product to a sale.
	DeleteProductInSale(ctx context.Context, id uuid.UUID) error
	// DeleteProductsBySaleID drops all products from sale.
	DeleteProductsBySaleID(ctx context.Context, saleID uuid.UUID) error
}
