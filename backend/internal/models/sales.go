package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Sale represents a completed purchase transaction in the system.
type Sale struct {
	ID        uuid.UUID       `gorm:"column:sale_id;type:uuid;primaryKey"`
	DateTime  time.Time       `gorm:"column:datetime;type:timestamp;not null"`
	Discount  decimal.Decimal `gorm:"column:discount;type:decimal(5,2)"`
	FullCost  decimal.Decimal `gorm:"column:full_cost;type:decimal(10,2);not null"`
	TotalCost decimal.Decimal `gorm:"column:total_cost;type:decimal(10,2);not null"`
}

// TableName returns table name from db.
func (Sale) TableName() string {
	return "sales"
}
