package gorm

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"supermarket/internal/models"
	base "supermarket/internal/repository/gorm"
	productType "supermarket/internal/repository/product_type"
)

type repository struct {
	*base.Repository[models.ProductType]

	db *gorm.DB
}

// New create gorm repository for product types.
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
	if err != nil {
		return nil, err
	}

	return pt, err
}
