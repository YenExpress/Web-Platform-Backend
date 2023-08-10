package routes

import (
	"github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/handlers"
	MW "github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/middlewares"

	"github.com/labstack/echo/v4"
)

func GroupRoutes(router echo.Echo) {

	group := router.Group("/auth/developer")

	group.POST("/login", func(c echo.Context) error {
		return handlers.Handler.HandleNativeLogin(c, MW.LoginDTOValidator)
	})

	group.POST("/signup", func(c echo.Context) error {
		return handlers.Handler.HandleNativeSignUp(c, MW.RegisterDTOValidator)
	})
}
