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
//
// @Summary      Get product supplies
// @Description  Returns a list of product supplies filtered by product, supplier and date range.
// @Tags         product-supplies
// @Produce      json
// @Param        product_id query string false "Product ID" format(uuid)
// @Param        supplier_id query string false "Supplier ID" format(uuid)
// @Param        date_from query string false "Start date (RFC3339)"
// @Param        date_to query string false "End date (RFC3339)"
// @Success      200 {array} dto.ProductSupplyDTO
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /product-supplies [get]
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
//
// @Summary      Create product supply
// @Description  Creates a new product supply.
// @Tags         product-supplies
// @Accept       json
// @Produce      json
// @Param        productSupply body models.ProductSupply true "Product supply"
// @Success      201 {object} models.ProductSupply
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /product-supplies [post]
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
//
// @Summary      Delete product supply
// @Description  Deletes a product supply by ID.
// @Tags         product-supplies
// @Produce      json
// @Param        id path string true "Product supply ID" format(uuid)
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      409 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /product-supplies/{id} [delete]
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
//
// @Summary      Update product supply
// @Description  Updates an existing product supply.
// @Tags         product-supplies
// @Accept       json
// @Produce      json
// @Param        id path string true "Product supply ID" format(uuid)
// @Param        productSupply body models.ProductSupply true "Product supply"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      409 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /product-supplies/{id} [put]
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
