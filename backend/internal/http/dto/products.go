package dto

import (
	"supermarket/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Product represents DTO for a concrete item in the catalog that can be sold or supplied.
type ProductDTO struct {
	ID           uuid.UUID       `json:"product_id"`
	Name         string          `json:"name"`
	Manufacturer string          `json:"manufacturer"`
	Weight       decimal.Decimal `json:"weight"`
	Volume       decimal.Decimal `json:"volume"`
	TypeName     string          `json:"type_name"`
}

// ToProductDTO converts model product to dto product.
func ToProductDTO(m *models.Product) ProductDTO {
	if m == nil {
		return ProductDTO{}
	}

	return ProductDTO{
		ID:           m.ID,
		Name:         m.Name,
		Manufacturer: m.Manufacturer,
		Weight:       m.Weight,
		Volume:       m.Volume,
		TypeName:     m.Type.Name,
	}
}
