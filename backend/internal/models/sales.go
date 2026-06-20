package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Sale represents a completed purchase transaction in the system.
type Sale struct {
	ID        uuid.UUID       `json:"sale_id" gorm:"column:sale_id;type:uuid;primaryKey"`
	DateTime  time.Time       `json:"datetime" gorm:"column:datetime;type:timestamp;not null"`
	Discount  decimal.Decimal `json:"discount" gorm:"column:discount;type:decimal(5,2)"`
	FullCost  decimal.Decimal `json:"full_cost" gorm:"column:full_cost;type:decimal(10,2);not null"`
	TotalCost decimal.Decimal `json:"total_cost" gorm:"column:total_cost;type:decimal(10,2);not null"`
}

// TableName returns table name from db.
func (Sale) TableName() string {
	return "sales"
}
