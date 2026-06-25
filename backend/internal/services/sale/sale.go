package sale

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"supermarket/internal/http/dto"
	"supermarket/internal/models"
	"supermarket/internal/repository/sale"
	productsale "supermarket/internal/services/product_sale"
)

// TODO: add transactions

// service provides business logic for sales.
type service struct {
	logger *slog.Logger

	salesR sale.Repository

	productSalesS productsale.Service
}

// New creates service for sales.
func New(
	logger *slog.Logger,
	salesR sale.Repository,
	productSalesS productsale.Service,
) Service {
	return &service{
		logger:        logger,
		salesR:        salesR,
		productSalesS: productSalesS,
	}
}

// GetSales returns list of sales filtered by date range.
func (s *service) GetSales(
	ctx context.Context,
	dateFrom *time.Time,
	dateTo *time.Time,
) ([]dto.SaleDTO, error) {
	const op = "services.sales.getSales"

	log := s.logger.With("op", op)

	salesList, err := s.salesR.GetByParams(ctx, dateFrom, dateTo)
	if err != nil {
		log.Error("failed to get sales",
			slog.Any("error", err),
			slog.Any("dateFrom", dateFrom),
			slog.Any("dateTo", dateTo),
		)

		return nil, err
	}

	result := make([]dto.SaleDTO, 0, len(salesList))
	for _, v := range salesList {
		result = append(result, dto.ToSaleDTO(&v))
	}

	return result, nil
}

// CreateSale creates sale header in the db.
func (s *service) CreateSale(ctx context.Context, sale *models.Sale) error {
	const op = "services.sales.createSale"

	log := s.logger.With("op", op)

	sale.ID = uuid.New()

	if err := s.salesR.Create(ctx, sale); err != nil {
		log.Error("failed to create sale",
			slog.Any("error", err),
			slog.Any("sale", *sale),
		)

		return err
	}

	return nil
}

// UpdateSale updates sale header in the db.
func (s *service) UpdateSale(ctx context.Context, sale *models.Sale) error {
	const op = "services.sales.updateSale"

	log := s.logger.With("op", op)

	if err := s.salesR.Update(ctx, sale); err != nil {
		log.Error("failed to update sale",
			slog.Any("error", err),
			slog.Any("sale", *sale),
		)

		return err
	}

	return nil
}

// DeleteSale deletes sale header in the db by id.
func (s *service) DeleteSale(ctx context.Context, id uuid.UUID) error {
	const op = "services.sales.deleteSale"

	log := s.logger.With("op", op).With("id", id)

	if err := s.productSalesS.DeleteProductsBySaleID(ctx, id); err != nil {
		log.Error("failed to delete products from sale", slog.Any("error", err))

		return err
	}

	if err := s.salesR.Delete(ctx, id); err != nil {
		log.Error("failed to delete sale", slog.Any("error", err))

		return err
	}

	return nil
}
