package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ignitedotdev Authentication Microservice!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
