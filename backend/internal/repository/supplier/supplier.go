package supplier

import (
	"context"

	"supermarket/internal/models"
	"supermarket/internal/repository"
)

// Repository is the interface for suppliers.
//
//go:generate mockery
type Repository interface {
	repository.Repository[models.Supplier]
	// GetByParams returns suppliers by search param.
	GetByParams(ctx context.Context, search string) ([]models.Supplier, error)
}
