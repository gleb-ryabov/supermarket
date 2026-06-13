package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Stock represents current inventory balance for a product.
type Stock struct {
	ID        uuid.UUID       `gorm:"column:stock_id;type:uuid;primaryKey"`
	ProductID *uuid.UUID      `gorm:"column:product_id;type:uuid"`
	Quantity  decimal.Decimal `gorm:"column:quantity;type:decimal(12,3)"`
}

// TableName returns table name from db.
func (Stock) TableName() string {
	return "stock"
}
