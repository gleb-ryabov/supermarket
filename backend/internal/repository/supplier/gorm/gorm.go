package gorm

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"supermarket/internal/models"
	base "supermarket/internal/repository/gorm"
	"supermarket/internal/repository/supplier"
)

type repository struct {
	*base.Repository[models.Supplier]

	db *gorm.DB
}

// New create gorm repository for suppliers.
func New(db *gorm.DB) supplier.Repository {
	return &repository{
		Repository: base.New[models.Supplier](db),
		db:         db,
	}
}

// GetByParams returns suppliers by search param.
func (r *repository) GetByParams(ctx context.Context, search string) ([]models.Supplier, error) {
	var p []models.Supplier

	q := r.db.WithContext(ctx)

	if search != "" {
		q = q.Where(
			"lower(name) like ? or lower(inn) like ? or lower(email) like ?",
			"%"+strings.ToLower(search)+"%",
			"%"+strings.ToLower(search)+"%",
			"%"+strings.ToLower(search)+"%",
		)
	}

	err := q.Find(&p).Error
	if err != nil {
		return nil, err
	}

	return p, nil
}
