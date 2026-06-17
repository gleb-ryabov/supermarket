package price

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	resp "supermarket/internal/lib/api/response"
	"supermarket/internal/models"
	"supermarket/internal/services/price"
)

// TODO: add swagger doc

const layoutDate = "2006-01-02" // layout for parse date.

// Handler is http request handler for price.
type Handler struct {
	logger *slog.Logger

	readTimeout  time.Duration
	writeTimeout time.Duration

	pricesS price.Service
}

// New creates handler for prices.
func New(
	logger *slog.Logger,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	pricesS price.Service,
) *Handler {
	return &Handler{
		logger:       logger,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		pricesS:      pricesS,
	}
}

// GetPrices handles HTTP GET request for retrieving prices list.
func (h *Handler) GetPrices(c *fiber.Ctx) error {
	const op = "http.handlers.price.getPrices"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.readTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	var typeID *uuid.UUID
	typeIDStr := c.Query("type_id")
	if typeIDStr != "" {
		v, err := uuid.Parse(typeIDStr)
		if err != nil {
			log.Error("failed to parse type_id to uuid", slog.Any("error", err), slog.String("typeId", typeIDStr))

			return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid type_id"))
		}
		typeID = &v
	}

	var dateFrom *time.Time
	dateFromStr := c.Query("date_from")
	if dateFromStr != "" {
		v, err := time.Parse(layoutDate, dateFromStr)
		if err != nil {
			log.Error(
				"failed to parse date_from to time.Time",
				slog.Any("error", err),
				slog.String("date", dateFromStr),
			)

			return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid date_from"))
		}
		dateFrom = &v
	}

	var dateTo *time.Time
	dateToStr := c.Query("date_to")
	if dateToStr != "" {
		v, err := time.Parse(layoutDate, dateToStr)
		if err != nil {
			log.Error("failed to parse date_to to time.Time", slog.Any("error", err), slog.String("date", dateToStr))

			return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid date_to"))
		}
		dateTo = &v
	}

	result, err := h.pricesS.GetPrices(ctx, typeID, dateFrom, dateTo)
	if err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed get prices"))
	}
	slog.Debug(
		"return sales",
		slog.Any("typeId", typeID),
		slog.Any("dateFrom", dateFrom),
		slog.Any("dateTo", dateTo),
	)

	return resp.Respond(c, fiber.StatusOK, result)
}

// CreatePrice handles HTTP POST request for create price in system.
func (h *Handler) CreatePrice(c *fiber.Ctx) error {
	const op = "http.handlers.prices.createPrice"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	var p models.Price
	if err := c.BodyParser(&p); err != nil {
		log.Error("failed to parse body", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid request body"))
	}

	if err := h.pricesS.CreatePrice(ctx, &p); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed create price"))
	}

	log.Debug("created price", slog.Any("price", p))

	return resp.Respond(c, fiber.StatusCreated, p)
}

// DeletePrice handles HTTP DELETE request for delete price in system.
func (h *Handler) DeletePrice(c *fiber.Ctx) error {
	op := "http.handlers.prices.deletePrice"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse id to uuid", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid id"))
	}

	if err = h.pricesS.DeletePrice(ctx, id); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed delete price"))
	}

	log.Debug("deleted price", slog.Any("id", id))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}

// UpdatePrice handles HTTP PUT request for update price in system.
func (h *Handler) UpdatePrice(c *fiber.Ctx) error {
	op := "http.handlers.prices.updatePrice"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse id to uuid", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid id"))
	}

	var p models.Price
	if err = c.BodyParser(&p); err != nil {
		log.Error("failed to parse body", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid request body"))
	}

	p.ID = id

	if err = h.pricesS.UpdatePrice(ctx, &p); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed update price"))
	}

	log.Debug("updated price", slog.Any("id", id), slog.Any("price", p))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}
