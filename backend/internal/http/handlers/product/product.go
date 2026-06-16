package product

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	resp "supermarket/internal/lib/api/response"
	"supermarket/internal/models"
	"supermarket/internal/services/product"
)

// TODO: add swagger doc

// Handler is http request handler for products.
type Handler struct {
	logger *slog.Logger

	readTimeout  time.Duration
	writeTimeout time.Duration

	productsS product.Service
}

// New creates handler for product types.
func New(
	logger *slog.Logger,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	productsS product.Service,
) *Handler {
	return &Handler{
		logger:       logger,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		productsS:    productsS,
	}
}

// GetProducts handles HTTP GET request for retrieving products list.
func (h *Handler) GetProducts(c *fiber.Ctx) error {
	const op = "http.handlers.products.getProducts"

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
	searchParam := c.Query("search")

	result, err := h.productsS.GetProducts(ctx, searchParam, typeID)
	if err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed get products"))
	}

	slog.Debug("return products", slog.Any("typeId", typeID), slog.String("search", searchParam))

	return resp.Respond(c, fiber.StatusOK, result)
}

// CreateProduct handles HTTP POST request for create product in system.
func (h *Handler) CreateProduct(c *fiber.Ctx) error {
	const op = "http.handlers.products.createProduct"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	var p models.Product
	if err := c.BodyParser(&p); err != nil {
		log.Error("failed to parse body", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid request body"))
	}

	if err := h.productsS.CreateProduct(ctx, &p); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed create product"))
	}

	log.Debug("created product", slog.Any("product", p))

	return resp.Respond(c, fiber.StatusCreated, p)
}

// DeleteProduct handles HTTP DELETE request for delete product in system.
func (h *Handler) DeleteProduct(c *fiber.Ctx) error {
	op := "http.handlers.products.deleteProduct"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse id to uuid", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid id"))
	}

	if err = h.productsS.DeleteProduct(ctx, id); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed delete product"))
	}

	log.Debug("deleted product", slog.Any("id", id))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}

// UpdateProduct handles HTTP PUT request for update product in system.
func (h *Handler) UpdateProduct(c *fiber.Ctx) error {
	op := "http.handlers.products.updateProduct"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse id to uuid", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid id"))
	}

	var p models.Product
	if err = c.BodyParser(&p); err != nil {
		log.Error("failed to parse body", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid request body"))
	}

	p.ID = id

	if err = h.productsS.UpdateProduct(ctx, &p); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed update product"))
	}

	log.Debug("updated product", slog.Any("id", id), slog.Any("product", p))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}
