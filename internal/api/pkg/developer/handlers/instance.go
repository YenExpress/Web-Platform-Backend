package handlers

import (
	I "github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/interactors"
	R "github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/repository"
)

var (
	repository = R.NewDeveloperRepository()
	service    = I.NewDeveloperAuthService(repository)
	Handler    = NewDeveloperAuthHandler(service)
)
