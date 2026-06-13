package gorm

import (
	"gorm.io/gorm"

	"supermarket/internal/models"
	base "supermarket/internal/repository/gorm"
	"supermarket/internal/repository/product"
)

type repository struct {
	*base.Repository[models.Product]
}

// New create gorm repository for products.
func New(db *gorm.DB) product.Repository {
	return repository{
		Repository: base.New[models.Product](db),
	}
}
