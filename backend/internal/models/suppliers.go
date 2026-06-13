package models

import "github.com/google/uuid"

// Supplier represents a vendor or company that supplies products to the system.
type Supplier struct {
	ID    uuid.UUID `gorm:"column:supplier_id;type:uuid;primaryKey"`
	Name  string    `gorm:"column:name;size:150;not null"`
	INN   string    `gorm:"column:inn;size:12"`
	KPP   string    `gorm:"column:kpp;size:9"`
	OGRN  string    `gorm:"column:ogrn;size:13"`
	Phone string    `gorm:"column:phone;size:11"`
	Email string    `gorm:"column:email;size:50"`
}

// TableName returns table name from db.
func (Supplier) TableName() string {
	return "suppliers"
}
