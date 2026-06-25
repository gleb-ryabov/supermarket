package gorm

import (
	"context"
	"time"

	"gorm.io/gorm"

	"supermarket/internal/models"
	base "supermarket/internal/repository/gorm"
	"supermarket/internal/repository/sale"
)

type repository struct {
	*base.Repository[models.Sale]

	db *gorm.DB
}

// New create gorm repository for sales.
func New(db *gorm.DB) sale.Repository {
	return &repository{
		Repository: base.New[models.Sale](db),
		db:         db,
	}
}

// GetByParams returns sales by params.
func (r *repository) GetByParams(
	ctx context.Context,
	dateFrom *time.Time,
	dateTo *time.Time,
) ([]models.Sale, error) {
	var sales []models.Sale

	q := r.db.WithContext(ctx)

	if dateFrom != nil {
		q = q.Where("datetime >= ?", dateFrom)
	}

	if dateTo != nil {
		q = q.Where("datetime <= ?", dateTo)
	}

	err := q.Find(&sales).Error
	if err != nil {
		return nil, err
	}

	return sales, nil
}
