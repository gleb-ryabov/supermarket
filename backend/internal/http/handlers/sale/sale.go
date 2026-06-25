package sale

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"supermarket/internal/lib/api/request"
	resp "supermarket/internal/lib/api/response"
	"supermarket/internal/models"
	"supermarket/internal/services/sale"
)

// Handler is http request handler for sales.
type Handler struct {
	logger *slog.Logger

	readTimeout  time.Duration
	writeTimeout time.Duration

	salesS sale.Service
}

// New creates handler for sales.
func New(
	logger *slog.Logger,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	salesS sale.Service,
) *Handler {
	return &Handler{
		logger:       logger,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		salesS:       salesS,
	}
}

// GetSales handles HTTP GET request for retrieving sales list.
func (h *Handler) GetSales(c *fiber.Ctx) error {
	const op = "http.handlers.sales.getSales"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.readTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	dateFrom, err := request.ParseQueryToTime(c, "date_from", time.RFC3339)
	if err != nil {
		log.Error("failed to parse date_from", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid date_from"))
	}

	dateTo, err := request.ParseQueryToTime(c, "date_to", time.RFC3339)
	if err != nil {
		log.Error("failed to parse date_to", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid date_to"))
	}

	result, err := h.salesS.GetSales(ctx, dateFrom, dateTo)
	if err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed get sales"))
	}

	log.Debug("return sales",
		slog.Any("dateFrom", dateFrom),
		slog.Any("dateTo", dateTo),
	)

	return resp.Respond(c, fiber.StatusOK, result)
}

// CreateSale handles HTTP POST request for create sale header.
func (h *Handler) CreateSale(c *fiber.Ctx) error {
	const op = "http.handlers.sales.createSale"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	var s models.Sale
	if err := c.BodyParser(&s); err != nil {
		log.Error("failed to parse body", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid request body"))
	}

	if err := h.salesS.CreateSale(ctx, &s); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed create sale"))
	}

	log.Debug("created sale", slog.Any("sale", s))

	return resp.Respond(c, fiber.StatusCreated, s)
}

// UpdateSale handles HTTP PUT request for update sale header.
func (h *Handler) UpdateSale(c *fiber.Ctx) error {
	const op = "http.handlers.sales.updateSale"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse id", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid id"))
	}

	var s models.Sale
	if err = c.BodyParser(&s); err != nil {
		log.Error("failed to parse body", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid request body"))
	}

	s.ID = id

	if err = h.salesS.UpdateSale(ctx, &s); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed update sale"))
	}

	log.Debug("updated sale", slog.Any("id", id), slog.Any("sale", s))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}

// DeleteSale handles HTTP DELETE request for delete sale header.
func (h *Handler) DeleteSale(c *fiber.Ctx) error {
	const op = "http.handlers.sales.deleteSale"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse id", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid id"))
	}

	if err = h.salesS.DeleteSale(ctx, id); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed delete sale"))
	}

	log.Debug("deleted sale", slog.Any("id", id))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}
