package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"supermarket/internal/http/handlers/price"
	"supermarket/internal/http/handlers/product"
	producttype "supermarket/internal/http/handlers/product_type"
	"supermarket/internal/http/handlers/stock"
	"supermarket/internal/http/handlers/supplier"
)

// Router is settings for http routs.
type Router struct {
	app *fiber.App

	productType *producttype.Handler
	product     *product.Handler
	price       *price.Handler
	supplier    *supplier.Handler
	stock       *stock.Handler
}

// New creates http router.
func New(
	app *fiber.App,
	productType *producttype.Handler,
	product *product.Handler,
	price *price.Handler,
	supplier *supplier.Handler,
	stock *stock.Handler,
) *Router {
	return &Router{
		app:         app,
		productType: productType,
		product:     product,
		price:       price,
		supplier:    supplier,
		stock:       stock,
	}
}

// Setup - configure router.
func (r *Router) Setup() {
	r.app.Use(recover.New())
	r.app.Use(requestid.New())
	r.app.Use(logger.New())
	r.app.Use(cors.New())

	api := r.app.Group("api")

	productTypes := api.Group("product-types")
	productTypes.Get("/", r.productType.GetProductTypes)
	productTypes.Post("/", r.productType.CreateProductType)
	productTypes.Delete("/:id", r.productType.DeleteProductType)
	productTypes.Put("/:id", r.productType.UpdateProductType)

	products := api.Group("products")
	products.Get("/", r.product.GetProducts)
	products.Post("/", r.product.CreateProduct)
	products.Delete("/:id", r.product.DeleteProduct)
	products.Put("/:id", r.product.UpdateProduct)

	prices := api.Group("/prices")
	prices.Get("/", r.price.GetPrices)
	prices.Post("/", r.price.CreatePrice)
	prices.Delete("/:id", r.price.DeletePrice)
	prices.Put("/:id", r.price.UpdatePrice)

	suppliers := api.Group("/suppliers")
	suppliers.Get("/", r.supplier.GetSuppliers)
	suppliers.Post("/", r.supplier.CreateSupplier)
	suppliers.Delete("/:id", r.supplier.DeleteSupplier)
	suppliers.Put("/:id", r.supplier.UpdateSupplier)

	stock := api.Group("/stock")
	stock.Get("/", r.stock.GetStocks)

	productSupplies := r.app.Group("/product-supplies")
	_ = productSupplies

	sales := r.app.Group("/sales")
	_ = sales

	productSales := r.app.Group("/product-sales")
	_ = productSales

	cancellations := r.app.Group("/cancellations")
	_ = cancellations

	reports := r.app.Group("/reports")
	_ = reports
}
