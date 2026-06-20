package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ProductSale represents a sale item within a sale transaction.
type ProductSale struct {
	ID        uuid.UUID       `json:"product_sales_id" gorm:"column:product_sales_id;type:uuid;primaryKey"`
	SaleID    uuid.UUID       `json:"sale_id" gorm:"column:sale_id;type:uuid;not null"`
	ProductID uuid.UUID       `json:"product_id" gorm:"column:product_id;type:uuid;not null"`
	SalePrice decimal.Decimal `json:"sale_price" gorm:"column:sale_price;type:decimal(10,2)"`
	Quantity  decimal.Decimal `json:"quantity" gorm:"column:quantity;type:decimal(12,3)"`

	Product Product `gorm:"foreign_key:product_id;references:product_id"`
	Sale    Sale    `gorm:"foreign_key:sale_id;references:sale_id"`
}

// TableName returns table name from db.
func (ProductSale) TableName() string {
	return "product_sales"
}
