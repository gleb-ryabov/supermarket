package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"supermarket/internal/models"
)

// ProductSupplyDTO represents DTO for product supply transaction.
type ProductSupplyDTO struct {
	ID           uuid.UUID       `json:"supply_id"`
	ProductID    uuid.UUID       `json:"product_id"`
	Price        decimal.Decimal `json:"price"`
	Quantity     decimal.Decimal `json:"quantity"`
	DeliveryDate string          `json:"delivery_date"`
	ProductName  string          `json:"product_name"`
	SupplierName string          `json:"supplier_name"`
}

// ToProductSupplyDTO converts model to dto.
func ToProductSupplyDTO(m *models.ProductSupply) ProductSupplyDTO {
	if m == nil {
		return ProductSupplyDTO{}
	}

	var deliveryDate string
	if m.DeliveryDate != nil {
		deliveryDate = m.DeliveryDate.Format("02.01.2006")
	}

	return ProductSupplyDTO{
		ID:           m.ID,
		ProductID:    m.ProductID,
		Price:        m.Price,
		Quantity:     m.Quantity,
		DeliveryDate: deliveryDate,
		ProductName:  m.Product.Name,
		SupplierName: m.Supplier.Name,
	}
}
