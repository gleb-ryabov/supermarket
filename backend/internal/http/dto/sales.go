package dto

import (
	"supermarket/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// SaleDTO represents sale transaction for API.
type SaleDTO struct {
	ID        uuid.UUID       `json:"sale_id"`
	DateTime  string          `json:"datetime"`
	Discount  decimal.Decimal `json:"discount"`
	FullCost  decimal.Decimal `json:"full_cost"`
	TotalCost decimal.Decimal `json:"total_cost"`
}

// ToSaleDTO converts model sale to dto.
func ToSaleDTO(m *models.Sale) SaleDTO {
	if m == nil {
		return SaleDTO{}
	}

	return SaleDTO{
		ID:        m.ID,
		DateTime:  m.DateTime.Format("02.01.2006 15:04"),
		Discount:  m.Discount,
		FullCost:  m.FullCost,
		TotalCost: m.TotalCost,
	}
}
