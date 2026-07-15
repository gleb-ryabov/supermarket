package gorm

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"supermarket/internal/models"
	repo "supermarket/internal/repository"
	base "supermarket/internal/repository/gorm"
	productsale "supermarket/internal/repository/product_sale"
)

type repository struct {
	*base.Repository[models.ProductSale]

	db *gorm.DB
}

// New create gorm repository for product in the sale.
func New(db *gorm.DB) productsale.Repository {
	return &repository{
		Repository: base.New[models.ProductSale](db),
		db:         db,
	}
}

// GetByParams returns product sales by params.
func (r *repository) GetBySale(ctx context.Context, saleID uuid.UUID) ([]models.ProductSale, error) {
	var ps []models.ProductSale

	if err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Sale").
		Where("sale_id = ?", saleID).
		Find(&ps).Error; err != nil {
		return nil, err
	}

	return ps, nil
}

// DeleteBySale drops product by sale.
func (r *repository) DeleteBySale(ctx context.Context, saleID uuid.UUID) error {
	var ps models.ProductSale

	res := r.db.WithContext(ctx).
		Where("sale_id = ?", saleID).
		Delete(&ps)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return repo.ErrNotFound
	}

	return nil
}
