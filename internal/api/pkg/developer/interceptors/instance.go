package interceptors

import (
	dev_dto "github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/dtm"
	"github.com/ignitedotdev/auth-ms/internal/api/shared/dto"
	middle "github.com/ignitedotdev/auth-ms/internal/api/shared/middlewares"
)

var (
	LoginDTOValidator    = new(middle.RequestBodyValidator[dto.LoginBody])
	RegisterDTOValidator = new(middle.RequestBodyValidator[dev_dto.RegisterDeveloperBody])
)
