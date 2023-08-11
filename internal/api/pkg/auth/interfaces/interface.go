package interfaces

import (
	SDTO "yenexpress/internal/api/pkg/auth/dto"
	SMW "yenexpress/internal/api/pkg/shared/middlewares"

	"github.com/labstack/echo/v4"
)

// interface describing shared behaviour across auth use cases
type IAuthService interface {
	NativeLogin(email, password string) error
}

// interface describing shared behaviour across auth handlers
type IAuthHandler interface {
	HandleNativeLogin(c echo.Context,
		validator *SMW.RequestBodyValidator[SDTO.LoginBody]) error
}
