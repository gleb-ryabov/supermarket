package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Price represents product pricing for a specific time period.
type Price struct {
	ID         uuid.UUID       `json:"price_id" gorm:"column:price_id;type:uuid;primaryKey"`
	ProductID  uuid.UUID       `json:"product_id" gorm:"column:product_id;type:uuid;not null"`
	DateStart  time.Time       `json:"date_start" gorm:"column:date_start;type:datetime;not null"`
	DateEnd    *time.Time      `json:"date_end" gorm:"column:date_end;type:datetime"`
	Discount   decimal.Decimal `json:"discount" gorm:"column:discount;type:decimal(5,2)"`
	FullPrice  decimal.Decimal `json:"full_price" gorm:"column:full_price;type:decimal(10,2)"`
	TotalPrice decimal.Decimal `json:"total_price" gorm:"column:total_price;type:decimal(10,2)"`

	Product Product `gorm:"foreign_key:product_id;references:product_id"`
}

// TableName returns table name from db.
func (Price) TableName() string {
	return "prices"
}
