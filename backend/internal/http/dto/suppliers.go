package dto

import (
	"github.com/google/uuid"

	"supermarket/internal/models"
)

// SupplierDTO represents supplier for API.
type SupplierDTO struct {
	ID    uuid.UUID `json:"supplier_id"`
	Name  string    `json:"name"`
	INN   string    `json:"inn"`
	KPP   string    `json:"kpp"`
	OGRN  string    `json:"ogrn"`
	Phone string    `json:"phone"`
	Email string    `json:"email"`
}

// ToSupplierDTO converts model supplier to dto.
func ToSupplierDTO(m *models.Supplier) SupplierDTO {
	if m == nil {
		return SupplierDTO{}
	}

	return SupplierDTO{
		ID:    m.ID,
		Name:  m.Name,
		INN:   m.INN,
		KPP:   m.KPP,
		OGRN:  m.OGRN,
		Phone: m.Phone,
		Email: m.Email,
	}
}
