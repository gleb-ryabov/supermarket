package product

import (
	"supermarket/internal/models"
	"supermarket/internal/repository"
)

// Repository is the interface for products.
type Repository interface {
	repository.Repository[models.Product]
}
