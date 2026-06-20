package models

import "github.com/google/uuid"

// Supplier represents a vendor or company that supplies products to the system.
type Supplier struct {
	ID    uuid.UUID `json:"supplier_id" gorm:"column:supplier_id;type:uuid;primaryKey"`
	Name  string    `json:"name" gorm:"column:name;size:150;not null"`
	INN   string    `json:"inn" gorm:"column:inn;size:12"`
	KPP   string    `json:"kpp" gorm:"column:kpp;size:9"`
	OGRN  string    `json:"ogrn" gorm:"column:ogrn;size:13"`
	Phone string    `json:"phone" gorm:"column:phone;size:11"`
	Email string    `json:"email" gorm:"column:email;size:50"`
}

// TableName returns table name from db.
func (Supplier) TableName() string {
	return "suppliers"
}
