package cancellation

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"supermarket/internal/lib/api/request"
	resp "supermarket/internal/lib/api/response"
	"supermarket/internal/models"
	"supermarket/internal/services/cancellation"
)

const layoutDate = "2006-01-02" // layout for parse date.

// Handler is http request handler for cancellation.
type Handler struct {
	logger *slog.Logger

	readTimeout  time.Duration
	writeTimeout time.Duration

	cancellationsS cancellation.Service
}

// New creates handler for cancellations.
func New(
	logger *slog.Logger,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	cancellationsS cancellation.Service,
) *Handler {
	return &Handler{
		logger:         logger,
		readTimeout:    readTimeout,
		writeTimeout:   writeTimeout,
		cancellationsS: cancellationsS,
	}
}

// GetCancellations handles HTTP GET request for retrieving cancellations list.
//
// @Summary      Get cancellations
// @Description  Returns a list of cancellations filtered by product and date range.
// @Tags         cancellations
// @Produce      json
// @Param        product_id query string false "Product ID" format(uuid)
// @Param        date_from  query string false "Start date (YYYY-MM-DD)"
// @Param        date_to    query string false "End date (YYYY-MM-DD)"
// @Success      200 {array} dto.CancellationDTO
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /cancellations [get].
func (h *Handler) GetCancellations(c *fiber.Ctx) error {
	const op = "http.handlers.cancellation.getCancellations"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.readTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	productID, err := request.ParseQueryToUUID(c, "product_id")
	if err != nil {
		log.Error("failed to parse product_id to uuid", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid product_id"))
	}

	dateFrom, err := request.ParseQueryToTime(c, "date_from", layoutDate)
	if err != nil {
		log.Error("failed to parse date_from to time.Time", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid date_from"))
	}

	dateTo, err := request.ParseQueryToTime(c, "date_to", layoutDate)
	if err != nil {
		log.Error("failed to parse date_to to time.Time", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid date_to"))
	}

	result, err := h.cancellationsS.GetCancellations(ctx, productID, dateFrom, dateTo)
	if err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed get cancellations"))
	}

	log.Debug(
		"return cancellations",
		slog.Any("productId", productID),
		slog.Any("dateFrom", dateFrom),
		slog.Any("dateTo", dateTo),
	)

	return resp.Respond(c, fiber.StatusOK, result)
}

// CreateCancellation handles HTTP POST request for create cancellation in system.
//
// @Summary      Create cancellations
// @Description  Creates a new cancellation.
// @Tags         cancellations
// @Accept       json
// @Produce      json
// @Param        cancellation body models.Cancellation true "Cancellation"
// @Success      201 {array} dto.CancellationDTO
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /cancellations [post].
func (h *Handler) CreateCancellation(c *fiber.Ctx) error {
	const op = "http.handlers.cancellations.createCancellation"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	var cancellation models.Cancellation
	if err := c.BodyParser(&cancellation); err != nil {
		log.Error("failed to parse body", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid request body"))
	}

	if err := h.cancellationsS.CreateCancellation(ctx, &cancellation); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed create cancellation"))
	}

	log.Debug("created cancellation", slog.Any("cancellation", cancellation))

	return resp.Respond(c, fiber.StatusCreated, cancellation)
}

// DeleteCancellation handles HTTP DELETE request for delete cancellation in system.
//
// @Summary      Delete cancellation
// @Description  Deletes a cancellation by ID.
// @Tags         cancellations
// @Accept       json
// @Produce      json
// @Param        id path string true "Cancellation ID" format(uuid)
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /cancellations/{id} [delete].
func (h *Handler) DeleteCancellation(c *fiber.Ctx) error {
	op := "http.handlers.cancellations.deleteCancellation"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse id to uuid", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid id"))
	}

	if err = h.cancellationsS.DeleteCancellation(ctx, id); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed delete cancellation"))
	}

	log.Debug("deleted cancellation", slog.Any("id", id))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}

// UpdateCancellation handles HTTP PUT request for update cancellation in system.
//
// @Summary      Update cancellation
// @Description  Updates an existing cancellation.
// @Tags         cancellations
// @Accept       json
// @Produce      json
// @Param        id path string true "Cancellation ID" format(uuid)
// @Param        cancellation body models.Cancellation true "Cancellation"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /cancellations/{id} [put].
func (h *Handler) UpdateCancellation(c *fiber.Ctx) error {
	op := "http.handlers.cancellations.updateCancellation"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse id to uuid", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid id"))
	}

	var cancellation models.Cancellation
	if err = c.BodyParser(&cancellation); err != nil {
		log.Error("failed to parse body", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid request body"))
	}

	cancellation.ID = id

	if err = h.cancellationsS.UpdateCancellation(ctx, &cancellation); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed update cancellation"))
	}

	log.Debug(
		"updated cancellation",
		slog.Any("id", id),
		slog.Any("cancellation", cancellation),
	)

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}
