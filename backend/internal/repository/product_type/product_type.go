package producttype

import (
	"context"

	"supermarket/internal/models"
	"supermarket/internal/repository"
)

// Repository is the interface for product types.
//
//go:generate mockery
type Repository interface {
	repository.Repository[models.ProductType]
	// GetByParams returns product types by params.
	GetByParams(ctx context.Context, name string, forAdult *bool) ([]models.ProductType, error)
}
