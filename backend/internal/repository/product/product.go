package product

import (
	"supermarket/internal/models"
	"supermarket/internal/repository"
)

type Repository interface {
	repository.Repository[models.Product]
}
