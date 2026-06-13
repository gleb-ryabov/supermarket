package sales

import (
	"supermarket/internal/models"
	"supermarket/internal/repository"
)

// Repository is the interface for sales.
type Repository interface {
	repository.Repository[models.Sale]
}
