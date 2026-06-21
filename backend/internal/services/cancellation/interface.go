package cancellation

import (
	"context"
	"time"

	"github.com/google/uuid"

	"supermarket/internal/http/dto"
	"supermarket/internal/models"
)

// Service provides business logic for cancellations.
type Service interface {
	// GetCancellations returns slice cancellations and error by product id, date from, date to.
	GetCancellations(
		ctx context.Context,
		productID *uuid.UUID,
		dateFrom *time.Time,
		dateTo *time.Time,
	) ([]dto.CancellationDTO, error)
	// CreateCancellation creates cancellation in the db.
	CreateCancellation(ctx context.Context, cancellation *models.Cancellation) error
	// DeleteCancellation deletes cancellation in the db by id.
	DeleteCancellation(ctx context.Context, id uuid.UUID) error
	// UpdateCancellation updates cancellation in the db.
	UpdateCancellation(ctx context.Context, cancellation *models.Cancellation) error
}
