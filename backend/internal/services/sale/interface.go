package sale

import (
	"context"
	"time"

	"github.com/google/uuid"

	"supermarket/internal/http/dto"
	"supermarket/internal/models"
)

// Service provides business logic for sales (header only).
type Service interface {
	// GetSales returns list of sales filtered by date range.
	GetSales(ctx context.Context, dateFrom *time.Time, dateTo *time.Time) ([]dto.SaleDTO, error)
	// CreateSale creates sale header in the db.
	CreateSale(ctx context.Context, sale *models.Sale) error
	// UpdateSale updates sale header in the db.
	UpdateSale(ctx context.Context, sale *models.Sale) error
	// DeleteSale deletes sale header in the db by id.
	DeleteSale(ctx context.Context, id uuid.UUID) error
}
