package gorm

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"supermarket/internal/models"
	base "supermarket/internal/repository/gorm"
	"supermarket/internal/repository/price"
)

type repository struct {
	*base.Repository[models.Price]

	db *gorm.DB
}

// New create gorm repository for prices.
func New(db *gorm.DB) price.Repository {
	return repository{
		Repository: base.New[models.Price](db),
		db:         db,
	}
}

// GetByParams returns prices by params.
func (r repository) GetByParams(
	ctx context.Context,
	typeID *uuid.UUID,
	dateFrom *time.Time,
	dateTo *time.Time,
) ([]models.Price, error) {
	var p []models.Price

	q := r.db.WithContext(ctx).Preload("Product")

	if typeID != nil {
		q = q.Where("type_id = ?", typeID)
	}

	if dateFrom != nil {
		q = q.Where("date_start >= ?", dateFrom)
	}

	if dateTo != nil {
		q = q.Where("date_end <= ? or date_end is null", dateTo)
	}

	err := q.Find(&p).Error
	if err != nil {
		return nil, err
	}

	return p, nil
}
