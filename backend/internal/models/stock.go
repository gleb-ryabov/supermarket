package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Stock represents current inventory balance for a product.
type Stock struct {
	ID        uuid.UUID       `json:"id" gorm:"column:stock_id;type:uuid;primaryKey"`
	ProductID *uuid.UUID      `json:"product_id" gorm:"column:product_id;type:uuid"`
	Quantity  decimal.Decimal `json:"quantity" gorm:"column:quantity;type:decimal(12,3)"`

	Product Product `gorm:"foreignKey:product_id;references:product_id"`
}

// TableName returns table name from db.
func (Stock) TableName() string {
	return "stock"
}
