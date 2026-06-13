package gorm

import (
	"supermarket/internal/models"
	base "supermarket/internal/repository/gorm"
	"supermarket/internal/repository/product"

	"gorm.io/gorm"
)

type repository struct {
	*base.Repository[models.Product]
}

func New(db *gorm.DB) product.Repository {
	return repository{
		Repository: base.New[models.Product](db),
	}
}
