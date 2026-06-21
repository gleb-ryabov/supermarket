package gorm

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"supermarket/internal/models"
	"supermarket/internal/repository/cancellation"
	base "supermarket/internal/repository/gorm"
)

type repository struct {
	*base.Repository[models.Cancellation]

	db *gorm.DB
}

// New create gorm repository for cancellations.
func New(db *gorm.DB) cancellation.Repository {
	return &repository{
		Repository: base.New[models.Cancellation](db),
		db:         db,
	}
}

// GetByParams returns cancellations by params.
func (r *repository) GetByParams(
	ctx context.Context,
	productID *uuid.UUID,
	dateFrom *time.Time,
	dateTo *time.Time,
) ([]models.Cancellation, error) {
	var c []models.Cancellation

	q := r.db.WithContext(ctx).Preload("Product")

	if productID != nil {
		q = q.Where("product_id = ?", productID)
	}

	if dateFrom != nil {
		q = q.Where("datetime >= ?", dateFrom)
	}

	if dateTo != nil {
		q = q.Where("datetime <= ?", dateTo)
	}

	err := q.Find(&c).Error
	if err != nil {
		return nil, err
	}

	return c, nil
}
