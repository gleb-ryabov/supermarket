package gorm

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository[T any] struct {
	db *gorm.DB
}

func New[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{
		db: db,
	}
}

func (r Repository[T]) GetAll(ctx context.Context) ([]*T, error) {
	var items []*T
	err := r.db.WithContext(ctx).Find(&items).Error

	return items, err
}

func (r Repository[T]) GetByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var item T
	err := r.db.WithContext(ctx).First(&item, id).Error

	return &item, err
}

func (r Repository[T]) Create(ctx context.Context, t *T) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r Repository[T]) Update(ctx context.Context, t *T) error {
	return r.db.WithContext(ctx).Model(t).
		Updates(t).Error
}

func (r Repository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	var item T
	return r.db.WithContext(ctx).Delete(item, id).Error
}
