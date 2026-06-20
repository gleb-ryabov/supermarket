package producttypes

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"supermarket/internal/lib/api/request"
	resp "supermarket/internal/lib/api/response"
	"supermarket/internal/models"
	producttype "supermarket/internal/services/product_type"
)

// TODO: add swagger doc

// Handler is http request handler for product types.
type Handler struct {
	logger *slog.Logger

	readTimeout  time.Duration
	writeTimeout time.Duration

	productTypesS producttype.Service
}

// New creates handler for product types.
func New(
	logger *slog.Logger,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	productTypeS producttype.Service,
) *Handler {
	return &Handler{
		logger:        logger,
		readTimeout:   readTimeout,
		writeTimeout:  writeTimeout,
		productTypesS: productTypeS,
	}
}

// GetProductTypes handles HTTP GET request for retrieving product types list.
func (h *Handler) GetProductTypes(c *fiber.Ctx) error {
	const op = "http.handlers.product_types.getProductTypes"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.readTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	forAdult, err := request.ParseQueryToBool(c, "for_adult")
	if err != nil {
		log.Error("failed to parse for_adult to bool", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid for_adult"))
	}
	searchParam := c.Query("search")

	result, err := h.productTypesS.GetProductTypes(ctx, searchParam, forAdult)
	if err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed get product types"))
	}

	log.Debug("return product types", slog.Any("for_adult", forAdult), slog.String("search", searchParam))

	return resp.Respond(c, fiber.StatusOK, result)
}

// CreateProductType handles HTTP POST request for create product type in system.
func (h *Handler) CreateProductType(c *fiber.Ctx) error {
	const op = "http.handlers.product_types.createProductType"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	var pt models.ProductType
	if err := c.BodyParser(&pt); err != nil {
		log.Error("failed to parse body", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid request body"))
	}

	if err := h.productTypesS.CreateProductType(ctx, &pt); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed create product type"))
	}

	log.Debug("created product type", slog.Any("productType", pt))

	return resp.Respond(c, fiber.StatusCreated, pt)
}

// DeleteProductType handles HTTP DELETE request for delete product type in system.
func (h *Handler) DeleteProductType(c *fiber.Ctx) error {
	op := "http.handlers.product_types.deleteProductType"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse id to uuid", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid id"))
	}

	if err = h.productTypesS.DeleteProductType(ctx, id); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed delete product type"))
	}

	log.Debug("deleted product type", slog.Any("id", id))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}

// UpdateProductType handles HTTP PUT request for update product type in system.
func (h *Handler) UpdateProductType(c *fiber.Ctx) error {
	op := "http.handlers.product_types.updateProductType"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse id to uuid", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid id"))
	}

	var pt models.ProductType
	if err = c.BodyParser(&pt); err != nil {
		log.Error("failed to parse body", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid request body"))
	}

	pt.ID = id

	if err = h.productTypesS.UpdateProductType(ctx, &pt); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed update product type"))
	}

	log.Debug("updated product type", slog.Any("id", id), slog.Any("productType", pt))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}
