package stock

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"

	"supermarket/internal/lib/api/request"
	resp "supermarket/internal/lib/api/response"
	"supermarket/internal/services/stock"
)

// TODO: add swagger doc

// Handler is http request handler for stock.
type Handler struct {
	logger *slog.Logger

	readTimeout  time.Duration
	writeTimeout time.Duration

	stockS stock.Service
}

// New creates handler for product types.
func New(
	logger *slog.Logger,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	stockS stock.Service,
) *Handler {
	return &Handler{
		logger:       logger,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		stockS:       stockS,
	}
}

// GetStocks handles HTTP GET request for retrieving stock list.
func (h *Handler) GetStocks(c *fiber.Ctx) error {
	const op = "http.handlers.stocks.getStocks"

	ctx, cancel := context.WithTimeout(c.UserContext(), h.readTimeout)
	defer cancel()

	log := h.logger.With("op", op).With("ip", c.IP())
	log.Debug("incoming request")

	productID, err := request.ParseQueryToUUID(c, "product_id")
	if err != nil {
		log.Error("failed to parse product_id to uuid", slog.Any("error", err))

		return resp.Respond(c, fiber.StatusBadRequest, resp.Error("invalid product_id"))
	}
	searchParam := c.Query("search")

	result, err := h.stockS.GetStocks(ctx, searchParam, productID)
	if err != nil {
		return resp.Respond(c, fiber.StatusInternalServerError, resp.Error("failed get stocks"))
	}

	slog.Debug("return stocks", slog.Any("productID", productID), slog.String("search", searchParam))

	return resp.Respond(c, fiber.StatusOK, result)
}
