package dto

import (
	"supermarket/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ProductSaleDTO represents product inside sale transaction.
type ProductSaleDTO struct {
	ID           uuid.UUID       `json:"product_sales_id"`
	SaleID       uuid.UUID       `json:"sale_id"`
	ProductID    uuid.UUID       `json:"product_id"`
	ProductName  string          `json:"product_name"`
	SalePrice    decimal.Decimal `json:"sale_price"`
	Quantity     decimal.Decimal `json:"quantity"`
	SaleDatetime string          `json:"sale_datetime"`
}

// ToProductSaleDTO converts model product sale to dto.
func ToProductSaleDTO(m *models.ProductSale) ProductSaleDTO {
	if m == nil {
		return ProductSaleDTO{}
	}

	return ProductSaleDTO{
		ID:           m.ID,
		SaleID:       m.SaleID,
		ProductID:    m.ProductID,
		ProductName:  m.Product.Name,
		SalePrice:    m.SalePrice,
		Quantity:     m.Quantity,
		SaleDatetime: m.Sale.DateTime.Format("02.01.2006 15:04"),
	}
}
