package routes

import (
	"github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/handlers"
	middle "github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/interceptors"

	"github.com/labstack/echo/v4"
)

func GroupRoutes(router echo.Echo) {

	group := router.Group("/auth/developer")

	group.POST("/login", func(c echo.Context) error {
		return handlers.Handler.HandleNativeLogin(c, middle.LoginDTOValidator)
	})

	group.POST("/signup", func(c echo.Context) error {
		return handlers.Handler.HandleNativeSignUp(c, middle.RegisterDTOValidator)
	})
}
