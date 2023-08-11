package prestart

import (
	"yenexpress/internal/api/database/entities"
	AR "yenexpress/internal/api/routes/authRoutes"

	"github.com/labstack/echo/v4"

	"yenexpress/internal/api/database/connectors"
)

// mount all routes associated with all handlers
func MountAllRoutes(r *echo.Echo) {

	AR.GroupPatientAuthRoutes(r)
}

// connect to database and set up all model instances for data persistence
func LoadDB() {
	connectors.ConnectDB(&entities.Patient{})
}
