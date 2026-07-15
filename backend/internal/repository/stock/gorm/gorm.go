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
	q := r.db.WithContext(ctx).
		Model(&models.Stock{}).
		Where("product_id = ?", productID)

	if count.IsNegative() {
		q = q.Where("quantity >= ?", count.Abs())
	}

	result := q.Update("quantity", gorm.Expr("quantity + ?", count))

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
		ProductID: &productID,
		Quantity:  decimal.NewFromInt(0),
	}

	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		FirstOrCreate(&s).Error

	if err != nil {
		return nil, err
	}

	return &s, nil
}

// UpdateStockByProductSale changes the quantity of a stock item by update product sale.
func (r *repository) UpdateStockByProductSale(
	ctx context.Context,
	productSaleID uuid.UUID,
	newQuantity decimal.Decimal,
) error {
	err := r.db.WithContext(ctx).Exec(`
			UPDATE stock s 
			SET quantity = s.quantity + (ps.quantity - ?)
			FROM product_sales ps
			WHERE ps.product_sales_id = ?
				AND ps.product_id = s.product_id 
		`, newQuantity, productSaleID,
	).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateStockOnDeleteSale changes the quantity of a stock item on drop sale.
func (r *repository) UpdateStockOnDeleteSale(ctx context.Context, saleID uuid.UUID) error {
	res := r.db.WithContext(ctx).Exec(`
		UPDATE stock s 
		SET quantity = s.quantity + ps.quantity
		FROM product_sales ps
		WHERE ps.sale_id = ?
			AND ps.product_id = s.product_id 
	`, saleID)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return repo.ErrNotFound
	}

	return nil
}

// UpdateStockByProductSupply changes the quantity of a stock item by update product supply.
func (r *repository) UpdateStockByProductSupply(
	ctx context.Context,
	productSupplyID uuid.UUID,
	newQuantity decimal.Decimal,
) error {
	res := r.db.WithContext(ctx).Exec(`
		UPDATE stock s 
		SET quantity = s.quantity - (ps.quantity - ?)
		FROM product_supplies ps
		WHERE ps.supply_id = ?
			AND ps.product_id = s.product_id 
	`, newQuantity, productSupplyID)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return repo.ErrNotFound
	}

	return nil
}

// UpdateStockByCancellation changes the quantity of a stock item by update cancellation.
func (r *repository) UpdateStockByCancellation(
	ctx context.Context,
	cancellationID uuid.UUID,
	newQuantity decimal.Decimal,
) error {
	res := r.db.WithContext(ctx).Exec(`
		UPDATE stock s 
		SET quantity = s.quantity + (c.quantity - ?)
		FROM cancellation c
		WHERE c.cancellation_id = ?
			AND c.product_id = s.product_id 
	`, newQuantity, cancellationID)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return repo.ErrNotFound
	}

	return nil
}
