package price

import (
	"context"
	"time"

	"github.com/google/uuid"

	"supermarket/internal/http/dto"
	"supermarket/internal/models"
)

// Service provides business logic for prices.
type Service interface {
	// GetPrices returns slice prices and error by params type id, date from, date to.
	GetPrices(ctx context.Context, typeID *uuid.UUID, dateFrom *time.Time, dateTo *time.Time) ([]dto.PriceDTO, error)
	// CreatePrice creates price in the db.
	CreatePrice(ctx context.Context, price *models.Price) error
	// DeletePrice deletes price in the db by id.
	DeletePrice(ctx context.Context, id uuid.UUID) error
	// UpdatePrice updates price in the db.
	UpdatePrice(ctx context.Context, price *models.Price) error
}
