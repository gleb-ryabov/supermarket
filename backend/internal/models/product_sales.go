package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ProductSale represents a sale item within a sale transaction.
type ProductSale struct {
	ID        uuid.UUID       `gorm:"column:product_sales_id;type:uuid;primaryKey"`
	SaleID    uuid.UUID       `gorm:"column:sale_id;type:uuid;not null"`
	ProductID uuid.UUID       `gorm:"column:product_id;type:uuid;not null"`
	SalePrice decimal.Decimal `gorm:"column:sale_price;type:decimal(10,2)"`
	Quantity  decimal.Decimal `gorm:"column:quantity;type:decimal(12,3)"`
}

// TableName returns table name from db.
func (ProductSale) TableName() string {
	return "product_sales"
}
