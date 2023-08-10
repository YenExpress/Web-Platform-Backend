package main

import (
	"net/http"

	"github.com/ignitedotdev/auth-ms/internal/api/config"
	"github.com/ignitedotdev/auth-ms/prestart"

	"github.com/labstack/echo/v4"
)

func init() {

	config.LoadConfig("/.env")

	prestart.LoadDB()

}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ignitedotdev Authentication Microservice!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
