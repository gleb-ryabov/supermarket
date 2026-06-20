package productsupply

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"supermarket/internal/models"
	productsupply "supermarket/internal/repository/product_supply"
	"supermarket/internal/services/stock"
)

// TODO: add transactions

// service provides business logic for product supplies.
type service struct {
	logger *slog.Logger

	stockS stock.Service

	productSupplyR productsupply.Repository
}

// New creates service for product supplies.
func New(
	logger *slog.Logger,
	productSupplyR productsupply.Repository,
	stockS stock.Service,
) Service {
	return &service{
		logger:         logger,
		stockS:         stockS,
		productSupplyR: productSupplyR,
	}
}

// GetProductSupplies returns slice product supplies and error by params product id, supplie id and period delivery.
func (s *service) GetProductSupplies(
	ctx context.Context,
	productID *uuid.UUID,
	supplierID *uuid.UUID,
	dateFrom *time.Time,
	dateTo *time.Time,
) ([]models.ProductSupply, error) {
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

	return pt, err
}

// CreateProductSupply creates product supply in the db.
func (s *service) CreateProductSupply(ctx context.Context, ps *models.ProductSupply) error {
	const op = "services.product_supply.createProductSupply"

	log := s.logger.With("op", op)

	ps.ID = uuid.New()

	if err := s.productSupplyR.Create(ctx, ps); err != nil {
		log.Error("failed to create product supply",
			slog.Any("error", err),
			slog.Any("productSupply", *ps),
		)

		return err
	}

	if err := s.stockS.SetCountStock(ctx, ps.ProductID, ps.Quantity); err != nil {
		log.Error("failed to set count stock",
			slog.Any("error", err),
			slog.Any("productSupply", *ps),
		)

		return err
	}

	return nil
}

// DeleteProductSupply deletes product supply in the db by id.
func (s *service) DeleteProductSupply(ctx context.Context, id uuid.UUID) error {
	const op = "services.product_supply.deleteProductSupply"

	log := s.logger.With("op", op)

	ps, err := s.productSupplyR.GetByID(ctx, id)
	if err != nil {
		log.Error("failed to get supply",
			slog.Any("error", err),
			slog.Any("id", id),
		)

		return err
	}

	if err = s.stockS.SetCountStock(ctx, ps.ProductID, ps.Quantity.Neg()); err != nil {
		log.Error("failed to set count stock",
			slog.Any("error", err),
			slog.Any("productSupply", *ps),
		)

		return err
	}

	if err = s.productSupplyR.Delete(ctx, id); err != nil {
		log.Error("failed to delete product supplies",
			slog.Any("error", err),
			slog.Any("id", id),
		)

		return err
	}

	return nil
}

// UpdateProductSupply deletes product supply in the db by id.
func (s *service) UpdateProductSupply(ctx context.Context, ps *models.ProductSupply) error {
	const op = "services.product_supply.updateProductSupply"

	log := s.logger.With("op", op)

	psOld, err := s.productSupplyR.GetByID(ctx, ps.ID)
	if err != nil {
		log.Error("failed to get old instance supply",
			slog.Any("error", err),
			slog.Any("id", ps.ID),
		)

		return err
	}

	shangedStock := ps.Quantity.Sub(psOld.Quantity)

	if err = s.stockS.SetCountStock(ctx, ps.ProductID, shangedStock); err != nil {
		log.Error("failed to create product supply",
			slog.Any("error", err),
			slog.Any("productSupply", *ps),
		)

		return err
	}

	if err = s.productSupplyR.Update(ctx, ps); err != nil {
		log.Error("failed to update product supply",
			slog.Any("error", err),
			slog.Any("productSupply", *ps),
		)

		return err
	}

	return nil
}
