package dto

import (
	"supermarket/internal/models"

	"github.com/google/uuid"
)

// ProductTypeDTO represents a DTO product category/type used for grouping products.
type ProductTypeDTO struct {
	ID       uuid.UUID `json:"type_id"`
	Name     string    `json:"name"`
	ForAdult bool      `json:"for_adult"`
}

// ToProductTypeDTO converts model product type to dto.
func ToProductTypeDTO(m *models.ProductType) ProductTypeDTO {
	if m == nil {
		return ProductTypeDTO{}
	}

	return ProductTypeDTO{
		ID:       m.ID,
		Name:     m.Name,
		ForAdult: m.ForAdult,
	}
}
