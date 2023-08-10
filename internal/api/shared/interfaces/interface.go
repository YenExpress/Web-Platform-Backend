package interfaces

import (
	SDTO "github.com/ignitedotdev/auth-ms/internal/api/shared/dto"
	SMW "github.com/ignitedotdev/auth-ms/internal/api/shared/middlewares"
	"github.com/labstack/echo/v4"
)

// interface describing common repository behaviours for different user models
type IRepository[userSchema interface{}] interface {
	SaveNew(user *userSchema) error
	GetByEmail(email string) (*userSchema, error)
	GetByID(ID interface{}) (*userSchema, error)
	Exists(ID interface{}) bool
}

// interface describing shared behaviour across use cases
type IService interface {
	NativeLogin(email, password string) error
}

// interface describing shared behaviour across handlers
type IHandler interface {
	HandleNativeLogin(c echo.Context,
		validator *SMW.RequestBodyValidator[SDTO.LoginBody]) error
}
