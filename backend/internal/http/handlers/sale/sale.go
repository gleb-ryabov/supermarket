package sale

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"supermarket/internal/lib/api/request"
	resp "supermarket/internal/lib/api/response"
	"supermarket/internal/models"
	"supermarket/internal/services"
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
//
// @Summary      Get sales
// @Description  Returns a list of sales filtered by date range.
// @Tags         sales
// @Produce      json
// @Param        date_from query string false "Start date (RFC3339)"
// @Param        date_to query string false "End date (RFC3339)"
// @Success      200 {array} dto.SaleDTO
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /sales [get].
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
//
// @Summary      Create sale
// @Description  Creates a new sale.
// @Tags         sales
// @Accept       json
// @Produce      json
// @Param        sale body models.Sale true "Sale"
// @Success      201 {object} models.Sale
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /sales [post].
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
//
// @Summary      Update sale
// @Description  Updates an existing sale.
// @Tags         sales
// @Accept       json
// @Produce      json
// @Param        id path string true "Sale ID" format(uuid)
// @Param        sale body models.Sale true "Sale"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /sales/{id} [put].
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
		switch {
		case errors.Is(err, services.ErrNotFound):
			return resp.Respond(c, fiber.StatusNotFound, resp.Error("sale not found"))
		default:
			return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed update sale"))
		}
	}

	log.Debug("updated sale", slog.Any("id", id), slog.Any("sale", s))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}

// DeleteSale handles HTTP DELETE request for delete sale header.
//
// @Summary      Delete sale
// @Description  Deletes a sale by ID.
// @Tags         sales
// @Produce      json
// @Param        id path string true "Sale ID" format(uuid)
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /sales/{id} [delete].
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
		switch {
		case errors.Is(err, services.ErrNotFound):
			return resp.Respond(c, fiber.StatusNotFound, resp.Error("sale not found"))
		default:
			return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed delete sale"))
		}
	}

	log.Debug("deleted sale", slog.Any("id", id))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}
