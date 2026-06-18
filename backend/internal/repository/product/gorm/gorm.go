package gorm

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"supermarket/internal/models"
	base "supermarket/internal/repository/gorm"
	"supermarket/internal/repository/product"
)

type repository struct {
	*base.Repository[models.Product]

	db *gorm.DB
}

// New create gorm repository for products.
func New(db *gorm.DB) product.Repository {
	return &repository{
		Repository: base.New[models.Product](db),
		db:         db,
	}
}

// GetByParams returns product types by params.
func (r *repository) GetByParams(ctx context.Context, name string, typeID *uuid.UUID) ([]models.Product, error) {
	var p []models.Product

	q := r.db.WithContext(ctx).Preload("Type")

	if typeID != nil {
		q = q.Where("type_id = ?", typeID)
	}

	if name != "" {
		q = q.Where("lower(name) like ?", "%"+strings.ToLower(name)+"%")
	}

	err := q.Find(&p).Error
	if err != nil {
		return nil, err
	}

	return p, nil
}
