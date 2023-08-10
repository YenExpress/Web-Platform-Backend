package handlers

import (
	"errors"
	"net/http"

	LDTO "github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/dto"
	SRC "github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/interactors"
	SH "github.com/ignitedotdev/auth-ms/internal/api/shared/handler"
	SMW "github.com/ignitedotdev/auth-ms/internal/api/shared/middlewares"

	EXC "github.com/ignitedotdev/auth-ms/internal/api/common/exceptions"

	"github.com/labstack/echo/v4"
)

// Handler object depends on use case object to accomplish tasks related to authentication
type DeveloperAuthHandler struct {
	SH.AuthHandler
	service *SRC.DeveloperAuthService
}

// create new developer authentication service providing handlers for different operations
func NewDeveloperAuthHandler(service *SRC.DeveloperAuthService) *DeveloperAuthHandler {
	return &DeveloperAuthHandler{service: service}
}

// handler function to handle developer sign up with provided information
// recieves input from request body and returns response object
func (handler *DeveloperAuthHandler) HandleNativeSignUp(c echo.Context,
	validator *SMW.RequestBodyValidator[LDTO.RegisterDeveloperBody]) error {
	credentials, err := validator.GetBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"Message": err.Error()})
		return err
	}

	if err = handler.service.NativeSignUp(
		credentials.FirstName, credentials.LastName,
		credentials.Email, credentials.Password,
	); err != nil {
		if errors.Is(err, EXC.UserExists) {
			c.JSON(http.StatusConflict, map[string]string{"Message": err.Error()})
			return err
		}
		c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
		return err
	}
	c.JSON(http.StatusCreated, map[string]string{"Message": "Developer Account Successfully Created"})
	return nil

}

// // handler function to handle developer login with credentials
// // recieves input credentials from request body and returns response object
// func (handler *DeveloperAuthHandler) HandleNativeLogin(c echo.Context,
// 	validator *SMW.RequestBodyValidator[SDTO.LoginBody]) error {
// 	credentials, err := validator.GetBody(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, map[string]string{"Message": err.Error()})
// 		return err
// 	}

// 	if err = handler.service.NativeLogin(credentials.Email, credentials.Password); err != nil {
// 		if errors.Is(err, EXC.UserDoesNotExist) {
// 			c.JSON(http.StatusNotFound, map[string]string{"Message": err.Error()})
// 			return err
// 		} else if errors.Is(err, EXC.InvalidPassword) {
// 			c.JSON(http.StatusUnauthorized, map[string]string{"Message": err.Error()})
// 			return err
// 		}
// 		c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
// 		return err
// 	}
// 	c.JSON(http.StatusCreated, map[string]string{"Message": "Admin Successfully Logged in"})
// 	return nil

// }
