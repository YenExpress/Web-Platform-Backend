package auth

import (
	"YenExpress/config"
	"YenExpress/service/patient/guard"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeWithAPIKey(c *gin.Context) bool {
	IPAddress, _ := config.GetIPAddress(c)
	apiKey, err := config.GetAPIKey(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
		return false
	}
	if guard.APIKeyLimiter.MaxOutFailure(apiKey, IPAddress) {
		c.JSON(http.StatusTooManyRequests, DefaultResponse{Message: "Failure To Validate API Key, Retry Later"})
		return false
	}
	if apiKeyIsValid, err := config.ValidateAPIKey(apiKey, config.PatientAPIKey); !apiKeyIsValid {
		guard.APIKeyLimiter.NoteFailure(apiKey, IPAddress)
		c.JSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
		return false
	}
	return true
}

func RateLimitLogin(c *gin.Context) (*LoginCredentials, bool) {
	var input *LoginCredentials
	IPAddress, _ := config.GetIPAddress(c)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, DefaultResponse{Message: err.Error()})
		return &LoginCredentials{}, false
	}
	if guard.LoginLimiter.MaxOutFailure(input.Email, IPAddress) {
		c.JSON(http.StatusTooManyRequests, DefaultResponse{Message: "Too Many Failed Login Attempts, Retry Later"})
		return &LoginCredentials{}, false
	}
	input.IPAddress = IPAddress
	return input, true
}

func RateLimitOTPValidation(c *gin.Context) (*valOTPCred, bool) {
	cred := &valOTPCred{}
	if err := cred.loadFromParams(c); err != nil {
		c.JSON(http.StatusBadRequest, DefaultResponse{Message: err.Error()})
		return &valOTPCred{}, false
	}
	if guard.EmailValidationLimiter.MaxOutFailure(cred.email, cred.ipAddress) {
		c.JSON(http.StatusTooManyRequests, DefaultResponse{Message: "Too Many Failed Login Attempts, Retry Later"})
		return &valOTPCred{}, false
	}
	return cred, true
}

func AuthorizeWithAccessToken(c *gin.Context) (*guard.Patient, bool) {
	token, err := config.GetBearerToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
		return &guard.Patient{}, false
	}
	user, err := guard.ValidateAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
		return &guard.Patient{}, false
	}
	return user, true
}

func AuthorizeWithRefreshToken(c *gin.Context) bool {
	token, err := config.GetBearerToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
		return false
	}
	err = guard.ValidateRefreshToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
		return false
	}
	return true
}
