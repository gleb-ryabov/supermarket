package price

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"supermarket/internal/http/dto"
	"supermarket/internal/models"
	"supermarket/internal/repository"
	"supermarket/internal/repository/price"
	"supermarket/internal/services"
)

// service provides business logic for prices.
type service struct {
	logger *slog.Logger

	priceR price.Repository
}

// New creates service for prices.
func New(
	logger *slog.Logger,
	priceR price.Repository,
) Service {
	return &service{
		logger: logger,
		priceR: priceR,
	}
}

// GetPrices returns slice prices and error by params type id, date from, date to.
func (s *service) GetPrices(
	ctx context.Context,
	typeID *uuid.UUID,
	dateFrom *time.Time,
	dateTo *time.Time,
) ([]dto.PriceDTO, error) {
	const op = "services.price.getPrices"

	log := s.logger.With("op", op)

	prices, err := s.priceR.GetByParams(ctx, typeID, dateFrom, dateTo)
	if err != nil {
		log.Error("failed to get products",
			slog.Any("error", err),
			slog.Any("typeId", typeID),
			slog.Any("dateFrom", dateFrom),
			slog.Any("dateTo", dateTo),
		)

		return nil, err
	}

	result := make([]dto.PriceDTO, 0, len(prices))
	for _, v := range prices {
		result = append(result, dto.ToPriceDTO(&v))
	}

	return result, nil
}

// CreatePrice creates price in the db.
func (s *service) CreatePrice(ctx context.Context, price *models.Price) error {
	const op = "services.price.createPrice"

	log := s.logger.With("op", op)

	price.ID = uuid.New()

	if err := s.priceR.Create(ctx, price); err != nil {
		log.Error("failed to create price",
			slog.Any("error", err),
			slog.Any("price", price),
		)

		return err
	}

	return nil
}

// DeletePrice deletes price in the db by id.
func (s *service) DeletePrice(ctx context.Context, id uuid.UUID) error {
	const op = "services.price.deletePrice"

	log := s.logger.With("op", op).
		With("id", id)

	if err := s.priceR.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			log.Error("price not found")

			return services.ErrNotFound
		default:
			log.Error("failed to delete price",
				slog.Any("error", err),
			)

			return err
		}
	}

	return nil
}

// UpdatePrice updates price in the db.
func (s *service) UpdatePrice(ctx context.Context, price *models.Price) error {
	const op = "services.price.updatePrice"

	log := s.logger.With("op", op).
		With("price", *price)

	if err := s.priceR.Update(ctx, price); err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			log.Error("price not found")

			return services.ErrNotFound
		default:
			log.Error("failed to update price",
				slog.Any("error", err),
			)

			return err
		}
	}

	return nil
}
