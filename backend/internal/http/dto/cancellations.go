package dto

import (
	"supermarket/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CancellationDTO represents product cancellation event for API.
type CancellationDTO struct {
	ID          uuid.UUID       `json:"cancellation_id"`
	DateTime    string          `json:"datetime"`
	ProductID   uuid.UUID       `json:"product_id"`
	ProductName string          `json:"product_name"`
	Quantity    decimal.Decimal `json:"quantity"`
}

// ToCancellationDTO converts model cancellation to dto.
func ToCancellationDTO(m *models.Cancellation) CancellationDTO {
	if m == nil {
		return CancellationDTO{}
	}

	return CancellationDTO{
		ID:          m.ID,
		DateTime:    m.DateTime.Format("02.01.2006 15:04"),
		ProductID:   m.ProductID,
		ProductName: m.Product.Name,
		Quantity:    m.Quantity,
	}
}
