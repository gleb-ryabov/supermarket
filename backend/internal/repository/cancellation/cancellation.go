package cancellation

import (
	"context"
	"time"

	"github.com/google/uuid"

	"supermarket/internal/models"
	"supermarket/internal/repository"
)

// Repository is the interface for cancellations.
//
//go:generate mockery
type Repository interface {
	repository.Repository[models.Cancellation]
	// GetByParams returns cancellation by params.
	GetByParams(
		ctx context.Context,
		productID *uuid.UUID,
		dateFrom *time.Time,
		dateTo *time.Time,
	) ([]models.Cancellation, error)
}
