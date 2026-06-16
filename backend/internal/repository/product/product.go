package product

import (
	"context"

	"github.com/google/uuid"

	"supermarket/internal/models"
	"supermarket/internal/repository"
)

// Repository is the interface for products.
type Repository interface {
	repository.Repository[models.Product]
	// GetByParams returns product types by params.
	GetByParams(ctx context.Context, name string, typeID *uuid.UUID) ([]models.Product, error)
}
