package handlers

import (
	"errors"
	"net/http"

	dev_dto "github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/dtm"
	"github.com/ignitedotdev/auth-ms/internal/api/pkg/developer/interactors"
	"github.com/ignitedotdev/auth-ms/internal/api/shared/dto"
	middle "github.com/ignitedotdev/auth-ms/internal/api/shared/middlewares"

	"github.com/ignitedotdev/auth-ms/internal/api/common/exceptions"

	"github.com/labstack/echo/v4"
)

// Handler object depends on use case object to accomplish tasks related to authentication
type DeveloperAuthHandler struct {
	service *interactors.DeveloperAuthService
}

// create new developer authentication service providing handlers for different operations
func NewDeveloperAuthHandler(service *interactors.DeveloperAuthService) *DeveloperAuthHandler {
	return &DeveloperAuthHandler{service: service}
}

// handler function to handle developer sign up with provided information
// recieves input from request body and returns response object
func (handler *DeveloperAuthHandler) HandleNativeSignUp(c echo.Context,
	validator *middle.RequestBodyValidator[dev_dto.RegisterDeveloperBody]) error {
	credentials, err := validator.GetBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"Message": err.Error()})
		return err
	}

	if err = handler.service.NativeSignUp(
		credentials.FirstName, credentials.LastName,
		credentials.Email, credentials.Password,
	); err != nil {
		if errors.Is(err, exceptions.UserExists) {
			c.JSON(http.StatusConflict, map[string]string{"Message": err.Error()})
			return err
		}
		c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
		return err
	}
	c.JSON(http.StatusCreated, map[string]string{"Message": "Developer Account Successfully Created"})
	return nil

}

// handler function to handle developer login with credentials
// recieves input credentials from request body and returns response object
func (handler *DeveloperAuthHandler) HandleNativeLogin(c echo.Context,
	validator *middle.RequestBodyValidator[dto.LoginBody]) error {
	credentials, err := validator.GetBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"Message": err.Error()})
		return err
	}

	if err = handler.service.NativeLogin(credentials.Email, credentials.Password); err != nil {
		if errors.Is(err, exceptions.UserDoesNotExist) {
			c.JSON(http.StatusNotFound, map[string]string{"Message": err.Error()})
			return err
		} else if errors.Is(err, exceptions.InvalidPassword) {
			c.JSON(http.StatusUnauthorized, map[string]string{"Message": err.Error()})
			return err
		}
		c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
		return err
	}
	c.JSON(http.StatusCreated, map[string]string{"Message": "Admin Successfully Logged in"})
	return nil

}
