package prestart

import (
	"github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/entities/models"
	"github.com/ignitedotdev/internal/api/common/utils"
	dev_route "github.com/ignitedotdev/internal/api/pkg/developer/routes"

	"github.com/ignitedotdev/internal/database"

	"github.com/labstack/echo/v4"
)

// Get all routes from all handlers
var (
	Routes = [][]utils.RouterFunc{dev_route.Group}
)

// mount all routes associated with all handlers
func MountAllRoutes(r echo.Echo) {
	utils.LoadRoutes(r, Routes...)
}

// connect to database and set up all model instances for data persistence
func LoadDB() {
	database.ConnectDB(&models.Developer{})
}
