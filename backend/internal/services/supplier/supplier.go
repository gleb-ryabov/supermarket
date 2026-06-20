package supplier

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"supermarket/internal/http/dto"
	"supermarket/internal/models"
	"supermarket/internal/repository/supplier"
)

// service provides business logic for suppliers.
type service struct {
	logger *slog.Logger

	supplierR supplier.Repository
}

// New creates service for suppliers.
func New(
	logger *slog.Logger,
	supplierR supplier.Repository,
) Service {
	return &service{
		logger:    logger,
		supplierR: supplierR,
	}
}

// GetSuppliers returns slice suppliers by search param.
func (s *service) GetSuppliers(ctx context.Context, search string) ([]dto.SupplierDTO, error) {
	const op = "services.supplier.getSuppliers"

	log := s.logger.With("op", op)

	suppliers, err := s.supplierR.GetByParams(ctx, search)
	if err != nil {
		log.Error("failed to get suppliers",
			slog.Any("error", err),
			slog.String("search", search),
		)

		return nil, err
	}

	result := make([]dto.SupplierDTO, 0, len(suppliers))
	for _, v := range suppliers {
		result = append(result, dto.ToSupplierDTO(&v))
	}

	return result, nil
}

// CreateSupplier creates supplier in the db.
func (s *service) CreateSupplier(ctx context.Context, supplier *models.Supplier) error {
	const op = "services.supplier.createSupplier"

	log := s.logger.With("op", op)

	supplier.ID = uuid.New()

	if err := s.supplierR.Create(ctx, supplier); err != nil {
		log.Error("failed to create supplier",
			slog.Any("error", err),
			slog.Any("supplier", supplier),
		)

		return err
	}

	return nil
}

// DeleteSupplier deletes supplier in the db by id.
func (s *service) DeleteSupplier(ctx context.Context, id uuid.UUID) error {
	const op = "services.supplier.deleteSupplier"

	log := s.logger.With("op", op)

	if err := s.supplierR.Delete(ctx, id); err != nil {
		log.Error("failed to delete supplier",
			slog.Any("error", err),
			slog.Any("id", id),
		)

		return err
	}

	return nil
}

// UpdateSupplier updates supplier in the db.
func (s *service) UpdateSupplier(ctx context.Context, sp *models.Supplier) error {
	const op = "services.supplier.updateSupplier"

	log := s.logger.With("op", op)

	if err := s.supplierR.Update(ctx, sp); err != nil {
		log.Error("failed to update supplier",
			slog.Any("error", err),
			slog.Any("supplier", *sp),
		)

		return err
	}

	return nil
}
