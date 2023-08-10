package repository

import (
	E "github.com/ignitedotdev/auth-ms/internal/api/database/entities"
	R "github.com/ignitedotdev/auth-ms/internal/api/shared/repository"
)

type DeveloperRepository struct {
	R.BRepository[E.Developer]
}

func NewDeveloperRepository() *DeveloperRepository {
	return &DeveloperRepository{}
}
