package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Cancellation represents a cancellation event of a product in the system.
type Cancellation struct {
	ID        uuid.UUID       `gorm:"column:cancellation_id;type:uuid;primaryKey"`
	DateTime  time.Time       `gorm:"column:datetime;type:timestamp;not null"`
	ProductID uuid.UUID       `gorm:"column:product_id;type:uuid;not null"`
	Quantity  decimal.Decimal `gorm:"column:quantity;type:decimal(12,3)"`

	Product Product `gorm:"foreign_key:product_id;references:product_id"`
}

// TableName returns table name from db.
func (Cancellation) TableName() string {
	return "cancellation"
}
