package productsale

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"supermarket/internal/http/dto"
	"supermarket/internal/models"
	productsale "supermarket/internal/repository/product_sale"
	"supermarket/internal/services/stock"
)

// TODO: add transactions

// service provides business logic for product in sale.
type service struct {
	logger *slog.Logger

	productSaleR productsale.Repository

	stockS stock.Service
}

// New creates service for product in sale.
func New(
	logger *slog.Logger,
	productSaleR productsale.Repository,
	stockS stock.Service,
) Service {
	return &service{
		logger:       logger,
		productSaleR: productSaleR,
		stockS:       stockS,
	}
}

// GetProductsInSale returns slice products from sale.
func (s *service) GetProductsInSale(ctx context.Context, saleID uuid.UUID) ([]dto.ProductSaleDTO, error) {
	const op = "services.product_sale.getProductsInSale"

	log := s.logger.With("op", op)

	pt, err := s.productSaleR.GetBySale(ctx, saleID)
	if err != nil {
		log.Error("failed to get products by sale ID", slog.Any("error", err), slog.Any("saleID", saleID))

		return nil, err
	}

	result := make([]dto.ProductSaleDTO, 0, len(pt))
	for _, v := range pt {
		result = append(result, dto.ToProductSaleDTO(&v))
	}

	return result, err
}

// CreateProductInSale adds a product to a sale.
func (s *service) CreateProductInSale(ctx context.Context, productSale *models.ProductSale) error {
	const op = "services.product_sale.createProductInSale"

	log := s.logger.With("op", op)

	productSale.ID = uuid.New()

	if err := s.productSaleR.Create(ctx, productSale); err != nil {
		log.Error("failed to create product in sale",
			slog.Any("error", err),
			slog.Any("productSale", *productSale),
		)

		return err
	}

	if err := s.stockS.IncreaseStock(
		ctx,
		productSale.ProductID,
		productSale.Quantity.Neg(),
	); err != nil {
		log.Error("failed to set count stock",
			slog.Any("error", err),
			slog.Any("productSale", *productSale),
		)

		return err
	}

	return nil
}

// UpdateProductInSale updates a product to a sale.
func (s *service) UpdateProductInSale(ctx context.Context, productSale *models.ProductSale) error {
	const op = "services.product_sale.updateProductInSale"

	log := s.logger.With("op", op)

	if err := s.stockS.UpdateStockByProductSale(ctx, productSale.ID, productSale.Quantity); err != nil {
		log.Error("failed to set count stock", slog.Any("error", err),
			slog.Any("error", err),
			slog.Any("productSale", *productSale),
		)

		return err
	}

	if err := s.productSaleR.Update(ctx, productSale); err != nil {
		log.Error("failed to update product in sale",
			slog.Any("error", err),
			slog.Any("productSale", *productSale),
		)

		return err
	}

	return nil
}

// DeleteProductInSale drops a product to a sale.
func (s *service) DeleteProductInSale(ctx context.Context, id uuid.UUID) error {
	const op = "services.product_sale.deleteProductInSale"

	log := s.logger.With("op", op).With("id", id)

	if err := s.stockS.UpdateStockByProductSale(ctx, id, decimal.Zero); err != nil {
		log.Error("failed to set count stock", slog.Any("error", err))

		return err
	}

	if err := s.productSaleR.Delete(ctx, id); err != nil {
		log.Error("failed to delete product in sale", slog.Any("error", err))

		return err
	}

	return nil
}

// DeleteProductsBySaleID drops all products from sale.
func (s *service) DeleteProductsBySaleID(ctx context.Context, saleID uuid.UUID) error {
	const op = "services.product_sale.deleteProductsBySaleID"

	log := s.logger.With("op", op)

	if err := s.stockS.UpdateStockOnDeleteSale(ctx, saleID); err != nil {
		log.Error("failed to set count stock", slog.Any("error", err))

		return err
	}

	if err := s.productSaleR.DeleteBySale(ctx, saleID); err != nil {
		log.Error("failed to delete products by sale ID", slog.Any("error", err), slog.Any("saleID", saleID))

		return err
	}

	return nil
}
