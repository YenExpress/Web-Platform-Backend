package middlewares

import (
	"yenexpress/internal/api/pkg/auth/dto"
	M "yenexpress/internal/api/pkg/shared/middlewares"
)

var (
	LoginDTOValidator    = new(M.RequestBodyValidator[dto.LoginBody])
	RegisterDTOValidator = new(M.RequestBodyValidator[dto.RegisterPatientBody])
)
