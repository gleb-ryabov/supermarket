package productsale

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"supermarket/internal/lib/api/request"
	resp "supermarket/internal/lib/api/response"
	"supermarket/internal/models"
	productsale "supermarket/internal/services/product_sale"
)

// Handler is http request handler for product in sale.
type Handler struct {
	logger *slog.Logger

	readTimeout  time.Duration
	writeTimeout time.Duration

	productSaleS productsale.Service
}

// New creates handler for product in sale.
func New(
	logger *slog.Logger,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	productSaleS productsale.Service,
) *Handler {
	return &Handler{
		logger:       logger,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		productSaleS: productSaleS,
	}
}

// GetProductsInSale handles HTTP GET request for retrieving product in sale list.
//
// @Summary      Get products in sale
// @Description  Returns a list of products in a sale.
// @Tags         product-sales
// @Produce      json
// @Param        sale_id query string true "Sale ID" format(uuid)
// @Success      200 {array} dto.ProductSaleDTO
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /product-sales [get].
func (h *Handler) GetProductsInSale(c *fiber.Ctx) error {
	const op = "http.handlers.product_sale.getProductsInSale"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.readTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	saleID, err := request.ParseQueryToUUID(c, "sale_id")
	if err != nil {
		log.Error("failed to parse sale_id to bool", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid sale_id"))
	}

	result, err := h.productSaleS.GetProductsInSale(ctx, *saleID)
	if err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed get product in sale"))
	}

	log.Debug("return product types", slog.Any("saleID", saleID))

	return resp.Respond(c, fiber.StatusOK, result)
}

// CreateProductInSale handles HTTP POST request for create product in system.
//
// @Summary      Create product in sale
// @Description  Adds a product to a sale.
// @Tags         product-sales
// @Accept       json
// @Produce      json
// @Param        productSale body models.ProductSale true "Product sale"
// @Success      201 {object} models.ProductSale
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /product-sales [post].
func (h *Handler) CreateProductInSale(c *fiber.Ctx) error {
	const op = "http.handlers.product_sale.createProductInSale"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	var p models.ProductSale
	if err := c.BodyParser(&p); err != nil {
		log.Error("failed to parse body", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid request body"))
	}

	if err := h.productSaleS.CreateProductInSale(ctx, &p); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed create product in sale"))
	}

	log.Debug("created product in sale", slog.Any("productSale", p))

	return resp.Respond(c, fiber.StatusCreated, p)
}

// DeleteProductInSale handles HTTP DELETE request for delete product in system.
//
// @Summary      Delete product from sale
// @Description  Removes a product from a sale by ID.
// @Tags         product-sales
// @Produce      json
// @Param        id path string true "Product sale ID" format(uuid)
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /product-sales/{id} [delete].
func (h *Handler) DeleteProductInSale(c *fiber.Ctx) error {
	op := "http.handlers.product_sale.deleteProductInSale"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse id to uuid", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid id"))
	}

	if err = h.productSaleS.DeleteProductInSale(ctx, id); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed delete product in sale"))
	}

	log.Debug("deleted product in sale", slog.Any("id", id))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}

// UpdateProductInSale handles HTTP PUT request for update product in system.
//
// @Summary      Update product in sale
// @Description  Updates a product in a sale.
// @Tags         product-sales
// @Accept       json
// @Produce      json
// @Param        id path string true "Product sale ID" format(uuid)
// @Param        productSale body models.ProductSale true "Product sale"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /product-sales/{id} [put].
func (h *Handler) UpdateProductInSale(c *fiber.Ctx) error {
	op := "http.handlers.product_sale.updateProductInSale"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.writeTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Error("failed to parse id to uuid", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid id"))
	}

	var p models.ProductSale
	if err = c.BodyParser(&p); err != nil {
		log.Error("failed to parse body", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid request body"))
	}

	p.ID = id

	if err = h.productSaleS.UpdateProductInSale(ctx, &p); err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed update product"))
	}

	log.Debug("updated product in sale", slog.Any("id", id), slog.Any("productSale", p))

	return resp.Respond(c, fiber.StatusOK, resp.OK())
}
