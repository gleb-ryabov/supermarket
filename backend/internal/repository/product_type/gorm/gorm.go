package gorm

import (
	"context"
	"strings"
	"supermarket/internal/models"
	base "supermarket/internal/repository/gorm"
	productType "supermarket/internal/repository/product_type"

	"gorm.io/gorm"
)

type repository struct {
	*base.Repository[models.ProductType]
	db *gorm.DB
}

func New(db *gorm.DB) productType.Repository {
	return repository{
		Repository: base.New[models.ProductType](db),
		db:         db,
	}
}

func (r repository) GetByParams(ctx context.Context, name string, forAdult *bool) ([]models.ProductType, error) {
	var pt []models.ProductType

	q := r.db.WithContext(ctx)

	if forAdult != nil {
		q = q.Where("for_adult = ?", forAdult)
	}

	if name != "" {
		q = q.Where("lower(name) like ?", "%"+strings.ToLower(name)+"%")
	}

	err := q.Find(&pt).Error

	return pt, err
}
