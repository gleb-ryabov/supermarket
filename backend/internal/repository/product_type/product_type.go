package producttype

import (
	"context"
	"supermarket/internal/models"
	"supermarket/internal/repository"
)

type Repository interface {
	repository.Repository[models.ProductType]
	GetByParams(ctx context.Context, name string, forAdult *bool) ([]models.ProductType, error)
}
