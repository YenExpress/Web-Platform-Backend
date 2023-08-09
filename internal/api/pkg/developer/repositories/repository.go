package repositories

import (
	repo "github.com/ignitedotdev/auth-ms/internal/api/common/repositories"
	"github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/entities/models"
)

type DeveloperRepository struct {
	repo.BUserRepository[models.Developer]
}

func NewDeveloperRepository() *DeveloperRepository {
	return &DeveloperRepository{}
}
