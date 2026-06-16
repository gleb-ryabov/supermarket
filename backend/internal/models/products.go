package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Product represents a concrete item in the catalog that can be sold or supplied.
type Product struct {
	ID           uuid.UUID       `json:"product_id" gorm:"column:product_id;type:uuid;primaryKey"`
	TypeID       uuid.UUID       `json:"type_id" gorm:"column:type_id;type:uuid;not null"`
	Name         string          `json:"name" gorm:"column:name;size:100;not null"`
	Manufacturer string          `json:"manufacturer" gorm:"column:manufacturer;size:100"`
	Weight       decimal.Decimal `json:"weight" gorm:"column:weight;type:decimal(12,3)"`
	Volume       decimal.Decimal `json:"volume" gorm:"column:volume;type:decimal(12,3)"`

	Type ProductType `json:"type" gorm:"foreignKey:type_id;references:type_id"`
}

// TableName returns table name from db.
func (Product) TableName() string {
	return "products"
}
