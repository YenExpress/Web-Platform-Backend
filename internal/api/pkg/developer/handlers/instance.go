package handlers

import (
	"github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/interactors"
	repo "github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/repositories"
)

var (
	repository = repo.NewDeveloperRepository()
	service    = interactors.NewDeveloperAuthService(repository)
	Handler    = NewDeveloperAuthHandler(service)
)
