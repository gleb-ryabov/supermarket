package repository

import (
	"context"

	"github.com/google/uuid"
)

// Repository is the interface for base repository.
type Repository[T any] interface {
	// GetAll returns all items from db.
	GetAll(ctx context.Context) ([]*T, error)
	// GetByID returns item by id from db.
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)
	// Create makes item in the db.
	Create(ctx context.Context, p *T) error
	// Update replaces item in the db.
	Update(ctx context.Context, t *T) error
	// Delete drops item in the db.
	Delete(ctx context.Context, id uuid.UUID) error
}
