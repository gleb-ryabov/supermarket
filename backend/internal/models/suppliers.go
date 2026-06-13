package models

import "github.com/google/uuid"

type Supplier struct {
	ID    uuid.UUID `gorm:"column:supplier_id;type:uuid;primaryKey"`
	Name  string    `gorm:"column:name;size:150;not null"`
	INN   string    `gorm:"column:inn;size:12"`
	KPP   string    `gorm:"column:kpp;size:9"`
	OGRN  string    `gorm:"column:ogrn;size:13"`
	Phone string    `gorm:"column:phone;size:11"`
	Email string    `gorm:"column:email;size:50"`
}

func (Supplier) TableName() string {
	return "suppliers"
}
