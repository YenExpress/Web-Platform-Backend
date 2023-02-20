package patientware

import (
	"YenExpress/config"
	guard "YenExpress/guard/patientguard"
	"YenExpress/ratelimiter"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizePatientWithAPIKey(c *gin.Context) bool {
	IPAddress, _ := config.GetIPAddress(c)
	apiKey, err := config.GetAPIKey(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
		return false
	}
	if ratelimiter.PatientAPIKeyLimiter.MaxOutFailure(apiKey, IPAddress) {
		c.JSON(http.StatusTooManyRequests, DefaultResponse{Message: "Failure To Validate API Key, Retry Later"})
		return false
	}
	if apiKeyIsValid, err := config.ValidateAPIKey(apiKey, config.PatientAPIKey); !apiKeyIsValid {
		ratelimiter.PatientAPIKeyLimiter.NoteFailure(apiKey, IPAddress)
		c.JSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
		return false
	}
	return true
}

func RateLimitPatientLogin(c *gin.Context) (*LoginCredentials, bool) {
	var input *LoginCredentials
	IPAddress, _ := config.GetIPAddress(c)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, DefaultResponse{Message: err.Error()})
		return &LoginCredentials{}, false
	}
	if ratelimiter.PatientLoginLimiter.MaxOutFailure(input.Email, IPAddress) {
		c.JSON(http.StatusTooManyRequests, DefaultResponse{Message: "Too Many Failed Login Attempts, Retry Later"})
		return &LoginCredentials{}, false
	}
	input.IPAddress = IPAddress
	return input, true
}

func AuthorizePatientWithAccessToken(c *gin.Context) (*guard.Patient, bool) {
	token, err := config.GetBearerToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
		return &guard.Patient{}, false
	}
	user, err := guard.ValidatePatientAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
		return &guard.Patient{}, false
	}
	return user, true
}

func AuthorizePatientWithRefreshToken(c *gin.Context) bool {
	token, err := config.GetBearerToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
		return false
	}
	err = guard.ValidatePatientRefreshToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
		return false
	}
	return true
}
