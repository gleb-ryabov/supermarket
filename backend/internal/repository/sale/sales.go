package sale

import (
	"context"
	"time"

	"supermarket/internal/models"
	"supermarket/internal/repository"
)

// Repository is the interface for sales.
//
//go:generate mockery
type Repository interface {
	repository.Repository[models.Sale]
	// GetByParams returns sales by params.
	GetByParams(
		ctx context.Context,
		dateFrom *time.Time,
		dateTo *time.Time,
	) ([]models.Sale, error)
}
