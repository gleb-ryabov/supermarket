package gorm

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"supermarket/internal/models"
	base "supermarket/internal/repository/gorm"
	productsupply "supermarket/internal/repository/product_supply"
)

type repository struct {
	*base.Repository[models.ProductSupply]

	db *gorm.DB
}

// New create gorm repository for product supplies.
func New(db *gorm.DB) productsupply.Repository {
	return &repository{
		Repository: base.New[models.ProductSupply](db),
		db:         db,
	}
}

// GetByParams returns product supplies by params.
func (r *repository) GetByParams(
	ctx context.Context,
	productID *uuid.UUID,
	supplierID *uuid.UUID,
	dateFrom *time.Time,
	dateTo *time.Time,
) ([]models.ProductSupply, error) {
	var ps []models.ProductSupply

	q := r.db.WithContext(ctx)

	if productID != nil {
		q = q.Where("product_id = ?", productID)
	}

	if supplierID != nil {
		q = q.Where("supplier_id = ?", supplierID)
	}

	if dateFrom != nil {
		q = q.Where("delivery_date >= ?", dateFrom)
	}

	if dateTo != nil {
		q = q.Where("delivery_date <= ?", dateTo)
	}

	err := q.Find(&ps).Error
	if err != nil {
		return nil, err
	}

	return ps, nil
}
