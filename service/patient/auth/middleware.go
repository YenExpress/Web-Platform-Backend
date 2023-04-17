package auth

import (
	"YenExpress/config"
	"YenExpress/service/patient/guard"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeWithAPIKey() gin.HandlerFunc {

	return func(c *gin.Context) {
		IPAddress, _ := config.GetIPAddress(c)
		apiKey, err := config.GetAPIKey(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
			return
		}
		if !guard.APIKeyLimiter.AllowRequest(IPAddress) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, DefaultResponse{Message: "Failure To Validate API Key, Retry Later"})
			return
		}
		if apiKeyIsValid, err := config.ValidateAPIKey(apiKey, config.PatientAPIKey); !apiKeyIsValid {
			guard.APIKeyLimiter.UpdateRequest(IPAddress)
			c.AbortWithStatusJSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
			return
		}
		c.Next()
	}
}

func RateLimitLogin(c *gin.Context) (*LoginCredentials, bool) {
	var input *LoginCredentials
	IPAddress, _ := config.GetIPAddress(c)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, DefaultResponse{Message: err.Error()})
		return &LoginCredentials{}, false
	}
	if !guard.LoginLimiter.AllowRequest(input.Email, IPAddress) {
		c.JSON(http.StatusTooManyRequests, DefaultResponse{Message: "Too Many Failed Login Attempts, Retry Later"})
		return &LoginCredentials{}, false
	}
	input.IPAddress = IPAddress
	return input, true
}

func RateLimitOTPGeneration(c *gin.Context) (*valOTPCred, bool) {
	cred := &valOTPCred{}
	if err := cred.loadFromParams(c); err != nil {
		c.JSON(http.StatusBadRequest, DefaultResponse{Message: err.Error()})
		return &valOTPCred{}, false
	}
	if !guard.CreateOTPLimiter.AllowRequest(cred.email, cred.ipAddress) {
		c.JSON(http.StatusTooManyRequests, DefaultResponse{Message: "Too Many Failed Login Attempts, Retry Later"})
		return &valOTPCred{}, false
	}
	return cred, true
}

func RateLimitOTPValidation(c *gin.Context) (*valOTPCred, bool) {
	cred := &valOTPCred{}
	if err := cred.loadFromParams(c); err != nil {
		c.JSON(http.StatusBadRequest, DefaultResponse{Message: err.Error()})
		return &valOTPCred{}, false
	}
	if guard.EmailValidationLimiter.AllowRequest(cred.email, cred.ipAddress) {
		c.JSON(http.StatusTooManyRequests, DefaultResponse{Message: "Too Many Failed Login Attempts, Retry Later"})
		return &valOTPCred{}, false
	}
	return cred, true
}

func AuthorizeWithAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := config.GetBearerToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
			return
		}
		_, err = guard.ValidateAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
			return
		}
		c.Next()
	}
}

func AuthorizeWithRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := config.GetBearerToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
			return
		}
		err = guard.ValidateRefreshToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
			return
		}
		c.Next()
	}
}
