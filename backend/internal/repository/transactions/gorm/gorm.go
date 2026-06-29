package gorm

import (
	"context"

	"gorm.io/gorm"

	cancellation "supermarket/internal/repository/cancellation/gorm"
	productsale "supermarket/internal/repository/product_sale/gorm"
	productsupply "supermarket/internal/repository/product_supply/gorm"
	sale "supermarket/internal/repository/sale/gorm"
	stock "supermarket/internal/repository/stock/gorm"
	"supermarket/internal/repository/transactions"
)

// unitOfWork is transactor for DB.
type unitOfWork struct {
	db *gorm.DB
}

// New create a transactor for DB.
func New(db *gorm.DB) transactions.UnitOfWork {
	return &unitOfWork{
		db: db,
	}
}

// Do execute transaction in DB.
func (u *unitOfWork) Do(
	ctx context.Context,
	fn func(ctx context.Context, repos transactions.Repositories) error,
) error {
	return u.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			repos := transactions.Repositories{
				Cancellation:  cancellation.New(tx),
				ProductSale:   productsale.New(tx),
				ProductSupply: productsupply.New(tx),
				Stock:         stock.New(tx),
				Sale:          sale.New(tx),
			}

			return fn(ctx, repos)
		})
}
