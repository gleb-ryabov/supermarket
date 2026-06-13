package producttypes

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

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
		logger:       logger,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,

		productTypesS: productTypeS,
	}
}

// GetProductTypes handles HTTP GET request for retrieving product types list.
func (h *Handler) GetProductTypes(c *fiber.Ctx) error {
	const op = "http.handlers.product_types.getProductTypes"

	ctx, cancel := context.WithTimeout(context.Background(), h.readTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incomming request")

	forAdultStr := c.Query("for_adult")
	var forAdult *bool
	if forAdultStr != "" {
		v, err := strconv.ParseBool(forAdultStr)
		if err != nil {
			log.Error("failed to parse for_adult to bool", slog.Any("error", err))

			return c.Status(fiber.StatusBadRequest).SendString("invalid for_adult")
		}

		forAdult = &v
	}
	searchParam := c.Query("search")

	result, err := h.productTypesS.GetProductTypes(ctx, searchParam, forAdult)
	if err != nil {
		log.Error("failed to get product types",
			slog.Any("error", err),
			slog.Any("for_adult", forAdult),
			slog.String("search", searchParam),
		)

		return c.Status(fiber.StatusInternalServerError).SendString("failed get product types")
	}

	log.Debug("return product types", slog.Any("for_adult", forAdult), slog.String("search", searchParam))

	return c.Status(fiber.StatusOK).JSON(result)
}

// CreateProductType handles HTTP POST request for create product in system.
func (h *Handler) CreateProductType(c *fiber.Ctx) error {
	const op = "http.handlers.product_types.createProductType"

	ctx, cancel := context.WithTimeout(context.Background(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incomming request")

	var pt models.ProductType
	if err := c.BodyParser(&pt); err != nil {
		log.Error("failed to parse body", slog.Any("body", c.Body()), slog.Any("error", err))

		return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
	}

	if err := h.productTypesS.CreateProductType(ctx, &pt); err != nil {
		log.Error("failed to create product type", slog.Any("productType", pt), slog.Any("error", err))

		return c.Status(fiber.StatusInternalServerError).SendString("failed create product type")
	}

	log.Debug("created product type", slog.Any("productType", pt))

	return c.Status(fiber.StatusCreated).JSON(pt)
}

// DeleteProductType handles HTTP DELETE request for delete product in system.
func (h *Handler) DeleteProductType(c *fiber.Ctx) error {
	op := "http.handlers.product_types.deleteProductType"

	ctx, cancel := context.WithTimeout(context.Background(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incomming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse id to uuid", slog.Any("error", err))

		return c.Status(fiber.StatusBadRequest).SendString("invalid id")
	}

	if err = h.productTypesS.DeleteProductType(ctx, id); err != nil {
		log.Error("failed to delete product type", slog.Any("id", id), slog.Any("error", err))

		return c.Status(fiber.StatusInternalServerError).SendString("failed delete product type")
	}

	log.Debug("deleted product type", slog.Any("id", id))

	return c.Status(fiber.StatusOK).SendString("success")
}

// UpdateProductType handles HTTP DELETE request for update product in system.
func (h *Handler) UpdateProductType(c *fiber.Ctx) error {
	op := "http.handlers.product_types.updateProductType"

	ctx, cancel := context.WithTimeout(context.Background(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incomming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse id to uuid", slog.Any("error", err))

		return c.Status(fiber.StatusBadRequest).SendString("invalid id")
	}

	var pt models.ProductType
	if err = c.BodyParser(&pt); err != nil {
		log.Error("failed to parse body", slog.Any("body", c.Body()), slog.Any("error", err))

		return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
	}

	pt.ID = id

	if err = h.productTypesS.UpdateProductType(ctx, &pt); err != nil {
		log.Error(
			"failed to update product type",
			slog.Any("id", id),
			slog.Any("productType", pt),
			slog.Any("error", err),
		)

		return c.Status(fiber.StatusInternalServerError).SendString("failed update product type")
	}

	log.Debug("updated product type", slog.Any("id", id), slog.Any("productType", pt))

	return c.Status(fiber.StatusOK).SendString("success")
}
