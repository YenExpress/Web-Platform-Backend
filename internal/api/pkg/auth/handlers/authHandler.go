package handlers

import (
	"errors"
	"net/http"

	"yenexpress/internal/api/pkg/auth/dto"
	INT "yenexpress/internal/api/pkg/auth/interfaces"
	SMW "yenexpress/internal/api/pkg/shared/middlewares"

	EXC "yenexpress/internal/api/pkg/auth/exceptions"

	"github.com/labstack/echo/v4"
)

// Handler object depends on use case object to accomplish tasks related to authentication
type AuthHandler struct {
	service INT.IAuthService
}

// handler function to handle developer login with credentials
// recieves input credentials from request body and returns response object
func (handler *AuthHandler) HandleNativeLogin(c echo.Context,
	validator *SMW.RequestBodyValidator[dto.LoginBody]) error {
	credentials, err := validator.GetBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"Message": err.Error()})
		return err
	}

	if err = handler.service.NativeLogin(credentials.Email, credentials.Password); err != nil {
		if errors.Is(err, EXC.UserDoesNotExist) {
			c.JSON(http.StatusNotFound, map[string]string{"Message": err.Error()})
			return err
		} else if errors.Is(err, EXC.InvalidPassword) {
			c.JSON(http.StatusUnauthorized, map[string]string{"Message": err.Error()})
			return err
		}
		c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
		return err
	}
	c.JSON(http.StatusCreated, map[string]string{"Message": "Admin Successfully Logged in"})
	return nil

}
