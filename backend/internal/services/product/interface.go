package product

import (
	"context"

	"github.com/google/uuid"

	"supermarket/internal/models"
)

// Service provides business logic for products.
type Service interface {
	// GetProducts returns slice products and error by params type name and type id.
	GetProducts(ctx context.Context, name string, typeID *uuid.UUID) ([]models.Product, error)
	// CreateProduct creates product in the db.
	CreateProduct(ctx context.Context, product *models.Product) error
	// DeleteProduct deletes product in the db by id.
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	// UpdateProduct updates product in the db.
	UpdateProduct(ctx context.Context, pt *models.Product) error
}
