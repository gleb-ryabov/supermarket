package models

import "github.com/google/uuid"

// ProductType represents a product category/type used for grouping products.
type ProductType struct {
	ID       uuid.UUID `gorm:"column:type_id;type:uuid;primaryKey"`
	Name     string    `gorm:"column:name;size:100"`
	ForAdult bool      `gorm:"column:for_adult"`
}

// TableName returns table name from db.
func (ProductType) TableName() string {
	return "product_types"
}
