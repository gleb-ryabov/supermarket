package gorm

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"supermarket/internal/models"
	base "supermarket/internal/repository/gorm"
	"supermarket/internal/repository/stock"
)

type repository struct {
	*base.Repository[models.Stock]

	db *gorm.DB
}

// New create gorm repository for stocks.
func New(db *gorm.DB) stock.Repository {
	return repository{
		Repository: base.New[models.Stock](db),
		db:         db,
	}
}

// GetByParams returns stocks by search param.
func (r repository) GetByParams(ctx context.Context, search string, productID *uuid.UUID) ([]models.Stock, error) {
	var s []models.Stock

	q := r.db.WithContext(ctx).Joins("Product").Preload("Product")

	if search != "" {
		q = q.Where("lower(\"Product\".name) like ?", "%"+strings.ToLower(search)+"%")
	}

	if productID != nil {
		q.Where("product_id = ?", productID)
	}

	err := q.Find(&s).Error
	if err != nil {
		return nil, err
	}

	return s, nil
}
