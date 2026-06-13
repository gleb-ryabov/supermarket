package router

import (
	producttype "supermarket/internal/http/handlers/product_type"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// Router is settings for http routs.
type Router struct {
	app *fiber.App

	productType *producttype.Handler
}

// New creates http router.
func New(
	app *fiber.App,
	productType *producttype.Handler,
) *Router {
	return &Router{
		app:         app,
		productType: productType,
	}
}

// Setup - configure router.
func (r *Router) Setup() {
	r.app.Use(recover.New())
	r.app.Use(requestid.New())
	r.app.Use(logger.New())
	r.app.Use(cors.New())

	productTypes := r.app.Group("api/product-types")
	productTypes.Get("/", r.productType.GetProductTypes)
	productTypes.Post("/", r.productType.CreateProductType)
	productTypes.Delete("/:id", r.productType.DeleteProductType)
	productTypes.Put("/:id", r.productType.UpdateProductType)

	products := r.app.Group("/products")
	_ = products

	prices := r.app.Group("/prices")
	_ = prices

	suppliers := r.app.Group("/suppliers")
	_ = suppliers

	productSupplies := r.app.Group("/product-supplies")
	_ = productSupplies

	stock := r.app.Group("/stock")
	_ = stock

	sales := r.app.Group("/sales")
	_ = sales

	productSales := r.app.Group("/product-sales")
	_ = productSales

	cancellations := r.app.Group("/cancellations")
	_ = cancellations

	reports := r.app.Group("/reports")
	_ = reports
}
