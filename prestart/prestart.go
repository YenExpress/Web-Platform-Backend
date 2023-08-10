package prestart

import (
	"github.com/ignitedotdev/auth-ms/internal/api/database/entities"
	dev_route "github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/routes"
	"github.com/labstack/echo/v4"

	"github.com/ignitedotdev/auth-ms/internal/api/database/connectors"
)

// mount all routes associated with all handlers
func MountAllRoutes(r echo.Echo) {
	dev_route.GroupRoutes(r)
}

// connect to database and set up all model instances for data persistence
func LoadDB() {
	connectors.ConnectDB(&entities.Developer{})
}
