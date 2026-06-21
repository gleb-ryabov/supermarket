package cancellation

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"supermarket/internal/http/dto"
	"supermarket/internal/models"
	"supermarket/internal/repository/cancellation"
	"supermarket/internal/services/stock"
)

// TODO: add transactions

// service provides business logic for cancellations.
type service struct {
	logger *slog.Logger

	cancellationR cancellation.Repository

	stockS stock.Service
}

// New creates service for cancellations.
func New(
	logger *slog.Logger,
	cancellationR cancellation.Repository,
	stockS stock.Service,
) Service {
	return &service{
		logger:        logger,
		cancellationR: cancellationR,
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

	if err := s.cancellationR.Create(ctx, cancellation); err != nil {
		log.Error("failed to create cancellation",
			slog.Any("error", err),
			slog.Any("cancellation", cancellation),
		)

		return err
	}

	if err := s.stockS.SetCountStock(ctx, cancellation.ProductID, cancellation.Quantity.Neg()); err != nil {
		log.Error("failed to set count stock",
			slog.Any("error", err),
			slog.Any("cancellation", *cancellation),
		)

		return err
	}

	return nil
}

// DeleteCancellation deletes cancellation in the db by id.
func (s *service) DeleteCancellation(ctx context.Context, id uuid.UUID) error {
	const op = "services.cancellation.deleteCancellation"

	log := s.logger.With("op", op)

	cancellation, err := s.cancellationR.GetByID(ctx, id)
	if err != nil {
		log.Error("failed to get cancellation",
			slog.Any("error", err),
			slog.Any("id", id),
		)

		return err
	}

	if err = s.stockS.SetCountStock(ctx, cancellation.ProductID, cancellation.Quantity); err != nil {
		log.Error("failed to set count stock",
			slog.Any("error", err),
			slog.Any("cancellation", *cancellation),
		)

		return err
	}

	if err = s.cancellationR.Delete(ctx, id); err != nil {
		log.Error("failed to delete cancellation",
			slog.Any("error", err),
			slog.Any("id", id),
		)

		return err
	}

	return nil
}

// UpdateCancellation updates cancellation in the db.
func (s *service) UpdateCancellation(ctx context.Context, cancellation *models.Cancellation) error {
	const op = "services.cancellation.updateCancellation"

	log := s.logger.With("op", op)

	cancellationOld, err := s.cancellationR.GetByID(ctx, cancellation.ID)
	if err != nil {
		log.Error("failed to get old instance cancellation",
			slog.Any("error", err),
			slog.Any("id", cancellation.ID),
		)

		return err
	}
	log.Debug("get id")

	shangedStock := cancellation.Quantity.Sub(cancellationOld.Quantity)
	if err = s.stockS.SetCountStock(ctx, cancellationOld.ProductID, shangedStock.Neg()); err != nil {
		log.Error("failed to set count stock",
			slog.Any("error", err),
			slog.Any("cancellation", *cancellation),
		)

		return err
	}

	if err = s.cancellationR.Update(ctx, cancellation); err != nil {
		log.Error("failed to update cancellation",
			slog.Any("error", err),
			slog.Any("cancellation", *cancellation),
		)

		return err
	}
	log.Debug(" update")

	return nil
}
