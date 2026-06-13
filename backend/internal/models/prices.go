package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Price represents product pricing for a specific time period.
type Price struct {
	ID         uuid.UUID       `gorm:"column:price_id;type:uuid;primaryKey"`
	ProductID  uuid.UUID       `gorm:"column:product_id;type:uuid;not null"`
	DateStart  time.Time       `gorm:"column:date_start;type:datetime;not null"`
	DateEnd    *time.Time      `gorm:"column:date_end;type:datetime"`
	Discount   decimal.Decimal `gorm:"column:discount;type:decimal(5,2)"`
	FullPrice  decimal.Decimal `gorm:"column:full_price;type:decimal(10,2)"`
	TotalPrice decimal.Decimal `gorm:"column:total_price;type:decimal(10,2)"`
}

// TableName returns table name from db.
func (Price) TableName() string {
	return "prices"
}
