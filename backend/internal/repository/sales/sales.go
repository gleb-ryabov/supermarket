package sales

import (
	"supermarket/internal/models"
	"supermarket/internal/repository"
)

type Repository interface {
	repository.Repository[models.Sale]
}
