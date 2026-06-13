package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Product represents a concrete item in the catalog that can be sold or supplied.
type Product struct {
	ID           uuid.UUID       `gorm:"column:product_id;type:uuid;primaryKey"`
	TypeID       uuid.UUID       `gorm:"column:type_id;type:uuid;not null"`
	Name         string          `gorm:"column:name;size:100;not null"`
	Manufacturer string          `gorm:"column:manufacturer;size:100"`
	Weight       decimal.Decimal `gorm:"column:weight;type:decimal(12,3)"`
	Volume       decimal.Decimal `gorm:"column:volume;type:decimal(12,3)"`
}

// TableName returns table name from db.
func (Product) TableName() string {
	return "products"
}
