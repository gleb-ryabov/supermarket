package gorm

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"supermarket/internal/models"
	repo "supermarket/internal/repository"
	base "supermarket/internal/repository/gorm"
	"supermarket/internal/repository/stock"
)

type repository struct {
	*base.Repository[models.Stock]

	db *gorm.DB
}

// New create gorm repository for stocks.
func New(db *gorm.DB) stock.Repository {
	return &repository{
		Repository: base.New[models.Stock](db),
		db:         db,
	}
}

// GetByParams returns stocks by search param.
func (r *repository) GetByParams(ctx context.Context, search string, productID *uuid.UUID) ([]models.Stock, error) {
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

// SetCountByProductID updates the quantity of a stock item by its ID.
func (r *repository) SetCountByProductID(ctx context.Context, productID uuid.UUID, count decimal.Decimal) error {
	result := r.db.WithContext(ctx).
		Model(&models.Stock{}).
		Where("product_id = ?", productID).
		Where("quantity >= ?", count).
		Update("quantity", gorm.Expr("quantity + ?", count))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return repo.ErrNotEnoughStock
	}

	return nil
}

// FirstOrCreateByProductID finds the stock by product ID, otherwise if not found creates a new stock.
func (r *repository) FirstOrCreateByProductID(ctx context.Context, productID uuid.UUID) (*models.Stock, error) {
	s := models.Stock{
		ID:        uuid.New(),
		ProductID: &productID,
		Quantity:  decimal.NewFromInt(0),
	}

	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		FirstOrCreate(&s).Error
	if err != nil {
		return nil, err
	}

	return &s, err
}
