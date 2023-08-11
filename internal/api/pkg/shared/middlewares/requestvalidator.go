package middlewares

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Request Body Validator object exposing validation middleware function and
// validated data object getter
type RequestBodyValidator[dto interface{}] struct {
}

// Validation Middleware function to validate the request body
func (val *RequestBodyValidator[dto]) Validate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqBody *dto
		if err := c.Bind(&reqBody); err != nil {
			// If there's an error while parsing JSON, return a bad request response
			return c.JSON(http.StatusBadRequest, map[string]string{"Message": "Invalid request body"})
		}

		// Add the parsed request body to the context for later use
		c.Set("validatedRequestBody", reqBody)
		// Call the next middleware/handler in the chain
		return next(c)
	}
}

// function to get request body as validated data
func (val *RequestBodyValidator[dto]) GetBody(c echo.Context) (*dto, error) {
	validated := c.Get("validatedRequestBody")
	if validated == nil {
		return nil, errors.New("Validated data not found in request context")
	}
	reqBody, ok := validated.(*dto)
	if !ok {
		return nil, errors.New("Failed to retrieve validated data")
	}
	return reqBody, nil
}
