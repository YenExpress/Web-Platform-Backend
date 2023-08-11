package handlers

import (
	"errors"
	"net/http"

	SRC "yenexpress/internal/api/pkg/auth/interactors"
	SMW "yenexpress/internal/api/pkg/shared/middlewares"
	R "yenexpress/internal/api/pkg/shared/repositories"

	"yenexpress/internal/api/pkg/auth/dto"
	EXC "yenexpress/internal/api/pkg/auth/exceptions"

	"github.com/labstack/echo/v4"
)

// Handler object depends on use case object to accomplish tasks related to authentication
type PatientAuthHandler struct {
	AuthHandler
	service *SRC.PatientAuthService
}

// create new developer authentication service providing handlers for different operations
func NewPatientAuthHandler(service *SRC.PatientAuthService) *PatientAuthHandler {
	return &PatientAuthHandler{service: service}
}

// handler function to handle developer sign up with provided information
// recieves input from request body and returns response object
func (handler *PatientAuthHandler) HandleNativeSignUp(c echo.Context,
	validator *SMW.RequestBodyValidator[dto.RegisterPatientBody]) error {
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

var (
	repository = R.NewPatientRepository()
	service    = SRC.NewPatientAuthService(repository)
	Handler    = NewPatientAuthHandler(service)
)
