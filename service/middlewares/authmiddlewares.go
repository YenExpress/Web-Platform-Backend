package middlewares

import (
	"YenExpress/config"
	"YenExpress/helper"
	"YenExpress/service/dto"
	"YenExpress/service/guard"
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIKeyAuthorization() gin.HandlerFunc {

	return func(c *gin.Context) {
		IPAddress, _ := helper.GetIPAddress(c)
		apiKey, err := guard.GetAPIKey(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.DefaultResponse{Message: err.Error()})
			return
		}
		if !guard.APIKeyLimiter.AllowRequest(IPAddress) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, dto.DefaultResponse{Message: "Failure To Validate API Key, Retry Later"})
			return
		}
		if apiKeyIsValid, err := guard.ValidateAPIKey(apiKey, config.APIKey); !apiKeyIsValid {
			guard.APIKeyLimiter.UpdateRequest(IPAddress)
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.DefaultResponse{Message: err.Error()})
			return
		}
		c.Next()
	}
}
