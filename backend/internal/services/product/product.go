package product

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"supermarket/internal/http/dto"
	"supermarket/internal/models"
	"supermarket/internal/repository/product"
)

// service provides business logic for products.
type service struct {
	logger *slog.Logger

	productR product.Repository
}

// New creates service for products.
func New(
	logger *slog.Logger,
	productR product.Repository,
) Service {
	return &service{
		logger:   logger,
		productR: productR,
	}
}

// GetProducts returns slice products and error by params type name and type id.
func (s *service) GetProducts(ctx context.Context, name string, typeID *uuid.UUID) ([]dto.ProductDTO, error) {
	const op = "services.product.getProducts"

	log := s.logger.With("op", op)

	products, err := s.productR.GetByParams(ctx, name, typeID)
	if err != nil {
		log.Error("failed to get products",
			slog.Any("error", err),
			slog.String("name", name),
			slog.Any("typeId", typeID),
		)

		return nil, err
	}

	result := make([]dto.ProductDTO, 0, len(products))
	for _, v := range products {
		result = append(result, dto.ToProductDTO(&v))
	}

	return result, nil
}

// CreateProduct creates product in the db.
func (s *service) CreateProduct(ctx context.Context, product *models.Product) error {
	const op = "services.product.getProducts"

	log := s.logger.With("op", op)

	product.ID = uuid.New()

	if err := s.productR.Create(ctx, product); err != nil {
		log.Error("failed to create product",
			slog.Any("error", err),
			slog.Any("product", product),
		)

		return err
	}

	return nil
}

// DeleteProduct deletes product in the db by id.
func (s *service) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	const op = "services.products.deleteProduct"

	log := s.logger.With("op", op)

	if err := s.productR.Delete(ctx, id); err != nil {
		log.Error("failed to delete product",
			slog.Any("error", err),
			slog.Any("id", id),
		)

		return err
	}

	return nil
}

// UpdateProduct deletes product in the db by id.
func (s *service) UpdateProduct(ctx context.Context, pt *models.Product) error {
	const op = "services.products.updateProduct"

	log := s.logger.With("op", op)

	if err := s.productR.Update(ctx, pt); err != nil {
		log.Error("failed to update product",
			slog.Any("error", err),
			slog.Any("product", *pt),
		)

		return err
	}

	return nil
}
