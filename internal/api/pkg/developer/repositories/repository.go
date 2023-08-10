package repositories

import (
	"github.com/ignitedotdev/auth-ms/internal/api/database/entities"
	repo "github.com/ignitedotdev/auth-ms/internal/api/shared/repositories"
)

type DeveloperRepository struct {
	repo.BUserRepository[entities.Developer]
}

func NewDeveloperRepository() *DeveloperRepository {
	return &DeveloperRepository{}
}
