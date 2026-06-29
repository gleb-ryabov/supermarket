package cancellation

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"supermarket/internal/http/dto"
	"supermarket/internal/models"
	"supermarket/internal/repository/cancellation"
	"supermarket/internal/repository/transactions"
	"supermarket/internal/services/stock"
)

// service provides business logic for cancellations.
type service struct {
	logger *slog.Logger

	cancellationR cancellation.Repository
	uow           transactions.UnitOfWork

	stockS stock.Service
}

// New creates service for cancellations.
func New(
	logger *slog.Logger,
	cancellationR cancellation.Repository,
	uow transactions.UnitOfWork,
	stockS stock.Service,
) Service {
	return &service{
		logger:        logger,
		cancellationR: cancellationR,
		uow:           uow,
		stockS:        stockS,
	}
}

// GetCancellations returns slice cancellations and error by product id, date from, date to.
func (s *service) GetCancellations(
	ctx context.Context,
	productID *uuid.UUID,
	dateFrom *time.Time,
	dateTo *time.Time,
) ([]dto.CancellationDTO, error) {
	const op = "services.cancellation.getCancellations"

	log := s.logger.With("op", op)

	cancellations, err := s.cancellationR.GetByParams(ctx, productID, dateFrom, dateTo)
	if err != nil {
		log.Error("failed to get cancellations",
			slog.Any("error", err),
			slog.Any("productId", productID),
			slog.Any("dateFrom", dateFrom),
			slog.Any("dateTo", dateTo),
		)

		return nil, err
	}

	result := make([]dto.CancellationDTO, 0, len(cancellations))
	for _, v := range cancellations {
		result = append(result, dto.ToCancellationDTO(&v))
	}

	return result, nil
}

// CreateCancellation creates cancellation in the db.
func (s *service) CreateCancellation(ctx context.Context, cancellation *models.Cancellation) error {
	const op = "services.cancellation.createCancellation"

	log := s.logger.With("op", op)

	cancellation.ID = uuid.New()

	return s.uow.Do(ctx, func(ctx context.Context, repos transactions.Repositories) error {
		if err := repos.Cancellation.Create(ctx, cancellation); err != nil {
			log.Error("failed to create cancellation",
				slog.Any("error", err),
				slog.Any("cancellation", cancellation),
			)

			return err
		}

		if err := stock.IncreaseStock(
			ctx,
			s.logger,
			repos.Stock,
			cancellation.ProductID,
			cancellation.Quantity.Neg(),
		); err != nil {
			log.Error("failed to set count stock",
				slog.Any("error", err),
				slog.Any("cancellation", *cancellation),
			)

			return err
		}

		return nil
	})
}

// DeleteCancellation deletes cancellation in the db by id.
func (s *service) DeleteCancellation(ctx context.Context, id uuid.UUID) error {
	const op = "services.cancellation.deleteCancellation"

	log := s.logger.With("op", op).
		With("id", id)

	return s.uow.Do(ctx, func(ctx context.Context, repos transactions.Repositories) error {
		if err := repos.Stock.UpdateStockByCancellation(ctx, id, decimal.Zero); err != nil {
			log.Error("failed to set count stock", slog.Any("error", err))

			return err
		}

		if err := repos.Cancellation.Delete(ctx, id); err != nil {
			log.Error("failed to delete cancellation", slog.Any("error", err))

			return err
		}

		return nil
	})
}

// UpdateCancellation updates cancellation in the db.
func (s *service) UpdateCancellation(ctx context.Context, cancellation *models.Cancellation) error {
	const op = "services.cancellation.updateCancellation"

	log := s.logger.With("op", op)

	return s.uow.Do(ctx, func(ctx context.Context, repos transactions.Repositories) error {
		if err := repos.Stock.UpdateStockByCancellation(ctx, cancellation.ID, cancellation.Quantity); err != nil {
			log.Error("failed to set count stock",
				slog.Any("error", err),
				slog.Any("cancellation", *cancellation),
			)

			return err
		}

		if err := repos.Cancellation.Update(ctx, cancellation); err != nil {
			log.Error("failed to update cancellation",
				slog.Any("error", err),
				slog.Any("cancellation", *cancellation),
			)

			return err
		}

		return nil
	})
}
