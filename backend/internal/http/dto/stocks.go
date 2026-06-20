package dto

import (
	"supermarket/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// StockDTO represents current inventory balance for a product.
type StockDTO struct {
	ID          uuid.UUID       `json:"stock_id"`
	Quantity    decimal.Decimal `json:"quantity"`
	ProductName string          `json:"product_name"`
}

// ToStockDTO converts model stock to dto.
func ToStockDTO(m *models.Stock) StockDTO {
	if m == nil {
		return StockDTO{}
	}

	return StockDTO{
		ID:          m.ID,
		Quantity:    m.Quantity,
		ProductName: m.Product.Name,
	}
}
