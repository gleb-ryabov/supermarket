package stock

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"supermarket/internal/models"
	"supermarket/internal/repository"
)

// Repository is the interface for stocks.
type Repository interface {
	repository.Repository[models.Stock]
	// GetByParams returns stocks by search param.
	GetByParams(ctx context.Context, search string, productID *uuid.UUID) ([]models.Stock, error)
	// SetCountByProductID updates the quantity of a stock item by product ID.
	SetCountByProductID(ctx context.Context, productID uuid.UUID, count decimal.Decimal) error
	// FirstOrCreateByProductID finds the stock by product ID, otherwise if not found creates a new stock.
	FirstOrCreateByProductID(ctx context.Context, productID uuid.UUID) (*models.Stock, error)
	// UpdateStockByProductSale changes the quantity of a stock item by update product sale.
	UpdateStockByProductSale(ctx context.Context, productSaleID uuid.UUID, newQuantity decimal.Decimal) error
	// UpdateStockOnDeleteSale changes the quantity of a stock item on drop sale.
	UpdateStockOnDeleteSale(ctx context.Context, saleID uuid.UUID) error
	// UpdateStockByProductSupply changes the quantity of a stock item by update product supply.
	UpdateStockByProductSupply(ctx context.Context, productSupplyID uuid.UUID, newQuantity decimal.Decimal) error
	// UpdateStockByCancellation changes the quantity of a stock item by update cancellation.
	UpdateStockByCancellation(ctx context.Context, cancellationID uuid.UUID, newQuantity decimal.Decimal) error
}
