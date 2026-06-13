package repository

import (
	"context"

	"github.com/google/uuid"
)

type Repository[T any] interface {
	GetAll(ctx context.Context) ([]*T, error)
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)
	Create(ctx context.Context, p *T) error
	Update(ctx context.Context, t *T) error
	Delete(ctx context.Context, id uuid.UUID) error
}
