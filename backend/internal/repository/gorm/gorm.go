package gorm

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"supermarket/internal/repository"
)

// Repository is base gorm repository with owner generic methods.
type Repository[T any] struct {
	db *gorm.DB
}

// New create gorm base repository.
func New[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{
		db: db,
	}
}

// GetAll returns all items from db.
func (r *Repository[T]) GetAll(ctx context.Context) ([]*T, error) {
	var items []*T
	err := r.db.WithContext(ctx).Find(&items).Error

	return items, err
}

// GetByID returns item by id from db.
func (r *Repository[T]) GetByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var item T
	err := r.db.WithContext(ctx).First(&item, id).Error

	return &item, err
}

// Create makes item in the db.
func (r *Repository[T]) Create(ctx context.Context, t *T) error {
	return r.db.WithContext(ctx).Create(t).Error
}

// Update replaces item in the db.
func (r *Repository[T]) Update(ctx context.Context, t *T) error {
	res := r.db.WithContext(ctx).
		Model(t).
		Select("*").
		Updates(t)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

// Delete drops item in the db.
func (r *Repository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	var item T

	res := r.db.WithContext(ctx).Delete(&item, id)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}
