package dto

import (
	"supermarket/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PriceDTO represents DTO for product pricing for a specific time period.
type PriceDTO struct {
	ID          uuid.UUID       `json:"price_id"`
	DateStart   string          `json:"date_start"`
	DateEnd     string          `json:"date_end"`
	Discount    decimal.Decimal `json:"discount"`
	FullPrice   decimal.Decimal `json:"full_price"`
	TotalPrice  decimal.Decimal `json:"total_price"`
	ProductName string          `json:"product_name"`
}

// ToPriceDTO converts model price to dto price.
func ToPriceDTO(m *models.Price) PriceDTO {
	if m == nil {
		return PriceDTO{}
	}

	return PriceDTO{
		ID:          m.ID,
		DateStart:   m.DateStart.Format("02.01.2006"),
		DateEnd:     m.DateEnd.Format("02.01.2006"),
		Discount:    m.Discount,
		FullPrice:   m.FullPrice,
		TotalPrice:  m.TotalPrice,
		ProductName: m.Product.Name,
	}
}
