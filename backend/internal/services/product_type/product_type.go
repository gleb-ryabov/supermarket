package producttype

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"supermarket/internal/http/dto"
	"supermarket/internal/models"
	producttype "supermarket/internal/repository/product_type"
)

// service provides business logic for product types.
type service struct {
	logger *slog.Logger

	productTypesR producttype.Repository
}

// New creates service for product types.
func New(
	logger *slog.Logger,
	productTypesR producttype.Repository,
) Service {
	return &service{
		logger:        logger,
		productTypesR: productTypesR,
	}
}

// GetProductTypes returns slice product types and error by params product name and for adult.
func (s *service) GetProductTypes(ctx context.Context, name string, forAdult *bool) ([]dto.ProductTypeDTO, error) {
	const op = "services.product_types.getProductTypes"

	log := s.logger.With("op", op)

	pt, err := s.productTypesR.GetByParams(ctx, name, forAdult)
	if err != nil {
		log.Error("failed to get product types",
			slog.Any("error", err),
			slog.String("name", name),
			slog.Any("for_adult", forAdult),
		)

		return nil, err
	}

	result := make([]dto.ProductTypeDTO, 0, len(pt))
	for _, v := range pt {
		result = append(result, dto.ToProductTypeDTO(&v))
	}

	return result, err
}

// CreateProductType creates product type in the db.
func (s *service) CreateProductType(ctx context.Context, pt *models.ProductType) error {
	const op = "services.product_types.createProductType"

	log := s.logger.With("op", op)

	pt.ID = uuid.New()

	if err := s.productTypesR.Create(ctx, pt); err != nil {
		log.Error("failed to create product type",
			slog.Any("error", err),
			slog.Any("productType", *pt),
		)

		return err
	}

	return nil
}

// DeleteProductType deletes product type in the db by id.
func (s *service) DeleteProductType(ctx context.Context, id uuid.UUID) error {
	const op = "services.product_types.deleteProductType"

	log := s.logger.With("op", op)

	if err := s.productTypesR.Delete(ctx, id); err != nil {
		log.Error("failed to delete product types",
			slog.Any("error", err),
			slog.Any("id", id),
		)

		return err
	}

	return nil
}

// UpdateProductType deletes product type in the db by id.
func (s *service) UpdateProductType(ctx context.Context, pt *models.ProductType) error {
	const op = "services.product_types.updateProductType"

	log := s.logger.With("op", op)

	if err := s.productTypesR.Update(ctx, pt); err != nil {
		log.Error("failed to update product type",
			slog.Any("error", err),
			slog.Any("productType", *pt),
		)

		return err
	}

	return nil
}
