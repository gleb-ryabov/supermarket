package transactions

import (
	"context"

	"supermarket/internal/repository/cancellation"
	productsale "supermarket/internal/repository/product_sale"
	productsupply "supermarket/internal/repository/product_supply"
	"supermarket/internal/repository/sale"
	"supermarket/internal/repository/stock"
)

// UnitOfWork is transactor for DB.
type UnitOfWork interface {
	// Do execute transaction in DB.
	Do(ctx context.Context, fn func(ctx context.Context, repos Repositories) error) error
}

// Repositories contains all repository from app.
type Repositories struct {
	Cancellation  cancellation.Repository
	ProductSale   productsale.Repository
	ProductSupply productsupply.Repository
	Stock         stock.Repository
	Sale          sale.Repository
}
