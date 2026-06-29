package supplier

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	resp "supermarket/internal/lib/api/response"
	"supermarket/internal/models"
	"supermarket/internal/services/supplier"
)

// Handler is http request handler for suppliers.
type Handler struct {
	logger *slog.Logger

	readTimeout  time.Duration
	writeTimeout time.Duration

	suppliersS supplier.Service
}

// New creates handler for suppliers.
func New(
	logger *slog.Logger,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	suppliersS supplier.Service,
) *Handler {
	return &Handler{
		logger:       logger,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		suppliersS:   suppliersS,
	}
}

// GetSuppliers handles HTTP GET request for retrieving suppliers list.
//
// @Summary      Get suppliers
// @Description  Returns a list of suppliers filtered by search string.
// @Tags         suppliers
// @Produce      json
// @Param        search query string false "Search string"
// @Success      200 {array} dto.SupplierDTO
// @Failure      500 {object} response.Response
// @Router       /suppliers [get]
func (h *Handler) GetSuppliers(c *fiber.Ctx) error {
	const op = "http.handlers.suppliers.getSuppliers"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.readTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	searchParam := c.Query("search")

	result, err := h.suppliersS.GetSuppliers(ctx, searchParam)
	if err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed get suppliers"))
	}

	slog.Debug("return suppliers", slog.String("search", searchParam))

	return resp.Respond(c, fiber.StatusOK, result)
}

// CreateSupplier handles HTTP POST request for create supplier in system.
//
// @Summary      Create supplier
// @Description  Creates a new supplier.
// @Tags         suppliers
// @Accept       json
// @Produce      json
// @Param        supplier body models.Supplier true "Supplier"
// @Success      201 {object} models.Supplier
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /suppliers [post]
func (h *Handler) CreateSupplier(c *fiber.Ctx) error {
	const op = "http.handlers.suppliers.createSupplier"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	var s models.Supplier
	if err := c.BodyParser(&s); err != nil {
		log.Error("failed to parse body", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid request body"))
	}

	if err := h.suppliersS.CreateSupplier(ctx, &s); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed create supplier"))
	}

	log.Debug("created supplier", slog.Any("supplier", s))

	return resp.Respond(c, fiber.StatusCreated, s)
}

// DeleteSupplier handles HTTP DELETE request for delete supplier in system.
//
// @Summary      Delete supplier
// @Description  Deletes a supplier by ID.
// @Tags         suppliers
// @Produce      json
// @Param        id path string true "Supplier ID" format(uuid)
// @Success      201 {object} models.Supplier
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /suppliers [delete]
func (h *Handler) DeleteSupplier(c *fiber.Ctx) error {
	const op = "http.handlers.suppliers.deleteSupplier"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse supplier_id to uuid", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid supplier_id"))
	}

	if err = h.suppliersS.DeleteSupplier(ctx, id); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed delete supplier"))
	}

	log.Debug("deleted supplier", slog.Any("id", id))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}

// UpdateSupplier handles HTTP PUT request for update supplier in system.
//
// @Summary      Update supplier
// @Description  Updates a supplier by ID.
// @Tags         suppliers
// @Accept       json
// @Produce      json
// @Param        id path string true "Supplier ID" format(uuid)
// @Param        supplier body models.Supplier true "Supplier"
// @Success      201 {object} models.Supplier
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /suppliers [put]
func (h *Handler) UpdateSupplier(c *fiber.Ctx) error {
	const op = "http.handlers.suppliers.updateSupplier"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse supplier_id to uuid", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid supplier_id"))
	}

	var s models.Supplier
	if err = c.BodyParser(&s); err != nil {
		log.Error("failed to parse body", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid request body"))
	}

	s.ID = id

	if err = h.suppliersS.UpdateSupplier(ctx, &s); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed update supplier"))
	}

	log.Debug("updated supplier", slog.Any("id", id), slog.Any("supplier", s))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}
