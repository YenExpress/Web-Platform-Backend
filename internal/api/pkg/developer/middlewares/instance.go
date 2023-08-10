package middlewares

import (
	LDTO "github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/dto"
	SDTO "github.com/ignitedotdev/auth-ms/internal/api/shared/dto"
	M "github.com/ignitedotdev/auth-ms/internal/api/shared/middlewares"
)

var (
	LoginDTOValidator    = new(M.RequestBodyValidator[SDTO.LoginBody])
	RegisterDTOValidator = new(M.RequestBodyValidator[LDTO.RegisterDeveloperBody])
)
