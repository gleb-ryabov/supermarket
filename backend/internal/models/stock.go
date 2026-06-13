package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Stock struct {
	ID        uuid.UUID       `gorm:"column:stock_id;type:uuid;primaryKey"`
	ProductID *uuid.UUID      `gorm:"column:product_id;type:uuid"`
	Quantity  decimal.Decimal `gorm:"column:quantity;type:decimal(12,3)"`
}

func (Stock) TableName() string {
	return "stock"
}
