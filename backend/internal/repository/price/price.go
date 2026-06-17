package price

import (
	"context"
	"time"

	"github.com/google/uuid"

	"supermarket/internal/models"
	"supermarket/internal/repository"
)

// Repository is the interface for products.
type Repository interface {
	repository.Repository[models.Price]
	// GetByParams returns prices by params.
	GetByParams(ctx context.Context, typeID *uuid.UUID, dateFrom *time.Time, dateTo *time.Time) ([]models.Price, error)
}
