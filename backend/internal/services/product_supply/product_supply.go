package productsupply

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"supermarket/internal/http/dto"
	"supermarket/internal/models"
	productsupply "supermarket/internal/repository/product_supply"
	"supermarket/internal/repository/transactions"
	"supermarket/internal/services/stock"
)

// service provides business logic for product supplies.
type service struct {
	logger *slog.Logger

	productSupplyR productsupply.Repository
	uow            transactions.UnitOfWork
}

// New creates service for product supplies.
func New(
	logger *slog.Logger,
	productSupplyR productsupply.Repository,
	uow transactions.UnitOfWork,
) Service {
	return &service{
		logger:         logger,
		productSupplyR: productSupplyR,
		uow:            uow,
	}
}

// GetProductSupplies returns slice product supplies and error by params product id, supplie id and period delivery.
func (s *service) GetProductSupplies(
	ctx context.Context,
	productID *uuid.UUID,
	supplierID *uuid.UUID,
	dateFrom *time.Time,
	dateTo *time.Time,
) ([]dto.ProductSupplyDTO, error) {
	const op = "services.product_supply.getProductSupplies"

	log := s.logger.With("op", op)

	pt, err := s.productSupplyR.GetByParams(ctx, productID, supplierID, dateFrom, dateTo)
	if err != nil {
		log.Error("failed to get product supplies",
			slog.Any("error", err),
			slog.Any("productId", productID),
			slog.Any("supplieId", supplierID),
			slog.Any("dateFrom", dateFrom),
			slog.Any("dateTo", dateTo),
		)

		return nil, err
	}

	result := make([]dto.ProductSupplyDTO, 0, len(pt))
	for _, v := range pt {
		result = append(result, dto.ToProductSupplyDTO(&v))
	}

	return result, err
}

// CreateProductSupply creates product supply in the db.
func (s *service) CreateProductSupply(ctx context.Context, ps *models.ProductSupply) error {
	const op = "services.product_supply.createProductSupply"

	log := s.logger.With("op", op)

	ps.ID = uuid.New()

	return s.uow.Do(ctx, func(ctx context.Context, repos transactions.Repositories) error {
		if err := repos.ProductSupply.Create(ctx, ps); err != nil {
			log.Error("failed to create product supply",
				slog.Any("error", err),
				slog.Any("productSupply", *ps),
			)

			return err
		}

		if err := stock.IncreaseStock(
			ctx,
			s.logger,
			repos.Stock,
			ps.ProductID,
			ps.Quantity,
		); err != nil {
			log.Error("failed to set count stock",
				slog.Any("error", err),
				slog.Any("productSupply", *ps),
			)

			return err
		}

		return nil
	})
}

// DeleteProductSupply deletes product supply in the db by id.
func (s *service) DeleteProductSupply(ctx context.Context, id uuid.UUID) error {
	const op = "services.product_supply.deleteProductSupply"

	log := s.logger.With("op", op).
		With("productSupplyID", id)

	return s.uow.Do(ctx, func(ctx context.Context, repos transactions.Repositories) error {
		if err := repos.Stock.UpdateStockByProductSupply(ctx, id, decimal.Zero); err != nil {
			log.Error("failed to set count stock", slog.Any("error", err))

			return err
		}

		if err := repos.ProductSupply.Delete(ctx, id); err != nil {
			log.Error("failed to delete product supplies", slog.Any("error", err))

			return err
		}

		return nil
	})
}

// UpdateProductSupply deletes product supply in the db by id.
func (s *service) UpdateProductSupply(ctx context.Context, ps *models.ProductSupply) error {
	const op = "services.product_supply.updateProductSupply"

	log := s.logger.With("op", op)

	return s.uow.Do(ctx, func(ctx context.Context, repos transactions.Repositories) error {
		if err := repos.Stock.UpdateStockByProductSupply(ctx, ps.ID, ps.Quantity); err != nil {
			log.Error("failed to set count stock",
				slog.Any("error", err),
				slog.Any("productSupply", *ps),
			)

			return err
		}

		if err := repos.ProductSupply.Update(ctx, ps); err != nil {
			log.Error("failed to update product supply",
				slog.Any("error", err),
				slog.Any("productSupply", *ps),
			)

			return err
		}

		return nil
	})
}
