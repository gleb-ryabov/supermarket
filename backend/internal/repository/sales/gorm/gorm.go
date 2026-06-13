package gorm

import (
	"supermarket/internal/models"
	base "supermarket/internal/repository/gorm"
	"supermarket/internal/repository/sales"

	"gorm.io/gorm"
)

type repository struct {
	*base.Repository[models.Sale]
}

func New(db *gorm.DB) sales.Repository {
	return &repository{
		Repository: base.New[models.Sale](db),
	}

}
