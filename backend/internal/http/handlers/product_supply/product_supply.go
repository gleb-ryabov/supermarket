package productsupply

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
	productsupply "supermarket/internal/services/product_supply"
)

// TODO: add swagger doc

// Handler is http request handler for product supplies.
type Handler struct {
	logger *slog.Logger

	readTimeout  time.Duration
	writeTimeout time.Duration

	productSuppliesS productsupply.Service
}

// New creates handler for product supplies.
func New(
	logger *slog.Logger,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	productSuppliesS productsupply.Service,
) *Handler {
	return &Handler{
		logger:           logger,
		readTimeout:      readTimeout,
		writeTimeout:     writeTimeout,
		productSuppliesS: productSuppliesS,
	}
}

// GetProductSupplies handles HTTP GET request for retrieving product supplies list.
func (h *Handler) GetProductSupplies(c *fiber.Ctx) error {
	const op = "http.handlers.product_supply.getProductSupplies"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.readTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	productID, err := request.ParseQueryToUUID(c, "product_id")
	if err != nil {
		log.Error("failed to parse productID id to UUID", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid product_id"))
	}

	supplierID, err := request.ParseQueryToUUID(c, "supplier_id")
	if err != nil {
		log.Error("failed to parse supplierID id to UUID", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid supplier_id"))
	}

	dateFrom, err := request.ParseQueryToTime(c, "date_from", time.RFC3339)
	if err != nil {
		log.Error("failed to parse dateFrom to time.Time", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid date_from"))
	}

	dateTo, err := request.ParseQueryToTime(c, "date_to", time.RFC3339)
	if err != nil {
		log.Error("failed to parse dateTo to time.Time", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid date_to"))
	}

	result, err := h.productSuppliesS.GetProductSupplies(ctx, productID, supplierID, dateFrom, dateTo)
	if err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed get product supplies"))
	}

	log.Debug(
		"return product supplies",
		slog.Any("productID", productID),
		slog.Any("supplierID", supplierID),
		slog.Any("dateFrom", dateFrom),
		slog.Any("dateTo", dateTo),
	)

	return resp.Respond(c, fiber.StatusOK, result)
}

// CreateProductSupply handles HTTP POST request for create product supply in system.
func (h *Handler) CreateProductSupply(c *fiber.Ctx) error {
	const op = "http.handlers.product_supply.createProductSupply"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	var ps models.ProductSupply
	if err := c.BodyParser(&ps); err != nil {
		log.Error("failed to parse body", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid request body"))
	}

	if err := h.productSuppliesS.CreateProductSupply(ctx, &ps); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed create product supply"))
	}

	log.Debug("created product supply", slog.Any("productSupply", ps))

	return resp.Respond(c, fiber.StatusCreated, ps)
}

// DeleteProductSupply handles HTTP DELETE request for delete product supply in system.
func (h *Handler) DeleteProductSupply(c *fiber.Ctx) error {
	op := "http.handlers.product_supply.deleteProductSupply"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse id to uuid", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid id"))
	}

	if err = h.productSuppliesS.DeleteProductSupply(ctx, id); err != nil {
		if errors.Is(err, services.ErrNotEnoughStock) {
			return resp.Respond(c, fiber.StatusConflict, resp.Error("not enough stock"))
		}

		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed delete product supply"))
	}

	log.Debug("deleted product supply", slog.Any("id", id))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}

// UpdateProductSupply handles HTTP PUT request for update product supply in system.
func (h *Handler) UpdateProductSupply(c *fiber.Ctx) error {
	op := "http.handlers.product_supply.updateProductSupply"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse id to uuid", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid id"))
	}

	var ps models.ProductSupply
	if err = c.BodyParser(&ps); err != nil {
		log.Error("failed to parse body", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid request body"))
	}

	ps.ID = id

	if err = h.productSuppliesS.UpdateProductSupply(ctx, &ps); err != nil {
		if errors.Is(err, services.ErrNotEnoughStock) {
			return resp.Respond(c, fiber.StatusConflict, resp.Error("not enough stock"))
		}

		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed update product supply"))
	}

	log.Debug("updated product supply", slog.Any("id", id), slog.Any("productSupply", ps))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}
