package gorm

import (
	"gorm.io/gorm"

	"supermarket/internal/models"
	base "supermarket/internal/repository/gorm"
	"supermarket/internal/repository/sales"
)

type repository struct {
	*base.Repository[models.Sale]
}

// New create gorm repository for sales.
func New(db *gorm.DB) sales.Repository {
	return &repository{
		Repository: base.New[models.Sale](db),
	}
}
