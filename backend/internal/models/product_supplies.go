package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ProductSupply represents a product supply transaction from a supplier.
type ProductSupply struct {
	ID           uuid.UUID       `json:"supply_id" gorm:"column:supply_id;type:uuid;primaryKey"`
	ProductID    uuid.UUID       `json:"product_id" gorm:"column:product_id;type:uuid;not null"`
	SupplierID   uuid.UUID       `json:"supplier_id" gorm:"column:supplier_id;type:uuid;not null"`
	Price        decimal.Decimal `json:"price" gorm:"column:price;type:decimal(10,2)"`
	Quantity     decimal.Decimal `json:"quantity" gorm:"column:quantity;type:decimal(12,3)"`
	DeliveryDate *time.Time      `json:"delivery_date" gorm:"column:delivery_date;type:date"`

	Product  Product  `gorm:"foreignKey:product_id;references:product_id"`
	Supplier Supplier `gorm:"foreignKey:supplier_id;references:supplier_id"`
}

// TableName returns table name from db.
func (ProductSupply) TableName() string {
	return "product_supplies"
}
