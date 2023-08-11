package main

import (
	"fmt"
	"net/http"

	"yenexpress/internal/api/config"
	"yenexpress/prestart"

	"github.com/labstack/echo/v4"
)

func init() {

	config.LoadConfig("/.env")

	prestart.LoadDB()

}

func main() {
	e := echo.New()
	// mount index route
	e.GET("/", func(c echo.Context) error {
		e.Logger.Info(fmt.Sprintf("Index Backend Service URL called by client with IP Address %v", c.RealIP()))
		return c.String(http.StatusOK, "The YenExpress Backend Service is Active!")
	})
	// mount all other api routes
	prestart.MountAllRoutes(e)
	e.Logger.Fatal(e.Start(":1323"))
}
