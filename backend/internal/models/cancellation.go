package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Cancellation struct {
	ID        uuid.UUID       `gorm:"column:cancellation_id;type:uuid;primaryKey"`
	DateTime  time.Time       `gorm:"column:datetime;type:timestamp;not null"`
	ProductID uuid.UUID       `gorm:"column:product_id;type:uuid;not null"`
	Quantity  decimal.Decimal `gorm:"column:quantity;type:decimal(12,3)"`
}

func (Cancellation) TableName() string {
	return "cancellation"
}
