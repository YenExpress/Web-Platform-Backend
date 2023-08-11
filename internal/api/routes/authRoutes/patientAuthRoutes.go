package authRoutes

import (
	"yenexpress/internal/api/pkg/auth/handlers"
	MW "yenexpress/internal/api/pkg/auth/middlewares"

	"github.com/labstack/echo/v4"
)

func GroupPatientAuthRoutes(router *echo.Echo) {

	group := router.Group("/auth/patient")

	group.POST("/login", func(c echo.Context) error {
		return handlers.Handler.HandleNativeLogin(c, MW.LoginDTOValidator)
	})

	group.POST("/signup", func(c echo.Context) error {
		return handlers.Handler.HandleNativeSignUp(c, MW.RegisterDTOValidator)
	})
}
