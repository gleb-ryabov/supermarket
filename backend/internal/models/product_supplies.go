package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ProductSupply represents a product supply transaction from a supplier.
type ProductSupply struct {
	ID           uuid.UUID       `gorm:"column:supply_id;type:uuid;primaryKey"`
	ProductID    uuid.UUID       `gorm:"column:product_id;type:uuid;not null"`
	SupplierID   uuid.UUID       `gorm:"column:supplier_id;type:uuid;not null"`
	Price        decimal.Decimal `gorm:"column:price;type:decimal(10,2)"`
	Quantity     decimal.Decimal `gorm:"column:quantity;type:decimal(12,3)"`
	DeliveryDate *time.Time      `gorm:"column:delivery_date;type:date"`
}

// TableName returns table name from db.
func (ProductSupply) TableName() string {
	return "product_supplies"
}
