package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"supermarket/internal/http/handlers/cancellation"
	"supermarket/internal/http/handlers/price"
	"supermarket/internal/http/handlers/product"
	productsale "supermarket/internal/http/handlers/product_sale"
	productsupply "supermarket/internal/http/handlers/product_supply"
	producttype "supermarket/internal/http/handlers/product_type"
	"supermarket/internal/http/handlers/sale"
	"supermarket/internal/http/handlers/stock"
	"supermarket/internal/http/handlers/supplier"
)

// Router is settings for http routs.
type Router struct {
	app *fiber.App

	productType   *producttype.Handler
	product       *product.Handler
	price         *price.Handler
	supplier      *supplier.Handler
	stock         *stock.Handler
	productSupply *productsupply.Handler
	cancellation  *cancellation.Handler
	productSale   *productsale.Handler
	sale          *sale.Handler
}

// New creates http router.
func New(
	app *fiber.App,
	productType *producttype.Handler,
	product *product.Handler,
	price *price.Handler,
	supplier *supplier.Handler,
	stock *stock.Handler,
	productSupply *productsupply.Handler,
	cancellation *cancellation.Handler,
	productSale *productsale.Handler,
	sale *sale.Handler,
) *Router {
	return &Router{
		app:           app,
		productType:   productType,
		product:       product,
		price:         price,
		supplier:      supplier,
		stock:         stock,
		productSupply: productSupply,
		cancellation:  cancellation,
		productSale:   productSale,
		sale:          sale,
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

	productSupplies := api.Group("/product-supplies")
	productSupplies.Get("/", r.productSupply.GetProductSupplies)
	productSupplies.Post("/", r.productSupply.CreateProductSupply)
	productSupplies.Delete("/:id", r.productSupply.DeleteProductSupply)
	productSupplies.Put("/:id", r.productSupply.UpdateProductSupply)

	productSales := api.Group("/product-sales")
	productSales.Get("/", r.productSale.GetProductsInSale)
	productSales.Post("/", r.productSale.CreateProductInSale)
	productSales.Delete("/:id", r.productSale.DeleteProductInSale)
	productSales.Put("/:id", r.productSale.UpdateProductInSale)

	sales := api.Group("/sales")
	sales.Get("/", r.sale.GetSales)
	sales.Post("/", r.sale.CreateSale)
	sales.Delete("/:id", r.sale.DeleteSale)
	sales.Put("/:id", r.sale.UpdateSale)

	cancellations := api.Group("/cancellations")
	cancellations.Get("/", r.cancellation.GetCancellations)
	cancellations.Post("/", r.cancellation.CreateCancellation)
	cancellations.Delete("/:id", r.cancellation.DeleteCancellation)
	cancellations.Put("/:id", r.cancellation.UpdateCancellation)

	reports := r.app.Group("/reports")
	_ = reports
}
