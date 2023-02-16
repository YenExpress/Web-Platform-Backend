package middleware

import (
	"YenExpress/guard"
	"YenExpress/ratelimiter"
	"YenExpress/util"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RateLimitPatientLogin(c *gin.Context) (body *LoginCredentials, status bool) {
	var input *LoginCredentials
	IPAddress, _ := util.GetIP(c)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, DefaultResponse{Message: err.Error()})
		return &LoginCredentials{}, false
	}
	limit := ratelimiter.PatientLoginLimiter.MaxOutFailure(input.Email, IPAddress)
	fmt.Println("Number of failed logins exceeded ? ", limit)
	if limit {
		c.JSON(http.StatusTooManyRequests, DefaultResponse{Message: "Too Many Failed Login Attempts, Retry Later"})
		return &LoginCredentials{}, false
	}
	input.IPAddress = IPAddress
	return input, true
}

func AuthorizePatientAccess(c *gin.Context) bool {
	status := true
	token := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
		_, err := guard.ValidatePatientAccessToken(token)
		if err != nil {
			c.JSON(http.StatusForbidden, DefaultResponse{Message: err.Error()})
			status = false
		}
	} else {
		c.JSON(http.StatusUnauthorized, DefaultResponse{Message: "Unauthorized!! Bearer Unknown"})
		status = false
	}
	return status
}

func AuthorizePatientRefresh(c *gin.Context) bool {
	status := true
	token := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
		err := guard.ValidatePatientRefreshToken(token)
		if err != nil {
			c.JSON(http.StatusForbidden, DefaultResponse{Message: err.Error()})
			status = false
		}
	} else {
		c.JSON(http.StatusUnauthorized, DefaultResponse{Message: "Unauthorized!! Bearer Unknown"})
		status = false
	}
	return status
}
