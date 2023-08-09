package middlewares

import (
	"github.com/ignitedotdev/auth-ms/internal/api/common/entities/dto"
	common_middle "github.com/ignitedotdev/auth-ms/internal/api/common/middlewares"
	dev_dto "github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/entities/dto"
)

var (
	LoginDTOValidator    = new(common_middle.RequestBodyValidator[dto.LoginCredentials])
	RegisterDTOValidator = new(common_middle.RequestBodyValidator[dev_dto.RegisterDeveloperDTO])
)
