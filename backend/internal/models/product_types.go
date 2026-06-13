package models

import "github.com/google/uuid"

type ProductType struct {
	ID       uuid.UUID `json:"type_id" gorm:"column:type_id;type:uuid;primaryKey"`
	Name     string    `json:"name" gorm:"column:name;size:100"`
	ForAdult bool      `json:"for_adult" gorm:"column:for_adult"`
}

func (ProductType) TableName() string {
	return "product_types"
}
