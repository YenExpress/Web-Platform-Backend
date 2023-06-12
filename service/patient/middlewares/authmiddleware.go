package middlewares

import (
	"YenExpress/helper"
	"YenExpress/service/dto"
	"YenExpress/service/guard"
	"YenExpress/service/patient/providers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RateLimitOTPGeneration(c *gin.Context) (*dto.OTPValidationCredentials, bool) {
	var input *dto.ConfirmEmail
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.DefaultResponse{Message: err.Error()})
		return &dto.OTPValidationCredentials{}, false
	}
	IPAddress, _ := helper.GetIPAddress(c)
	cred := &dto.OTPValidationCredentials{}
	cred.Email, cred.IPAddress = input.Email, IPAddress
	if !providers.CreateOTPLimiter.AllowRequest(cred.Email, cred.IPAddress) {
		c.JSON(http.StatusTooManyRequests, dto.DefaultResponse{Message: "Too Many OTP Generation Attempts, Retry Later"})
		return &dto.OTPValidationCredentials{}, false
	}
	return cred, true
}

func RateLimitOTPValidation(c *gin.Context) (*dto.OTPValidationCredentials, bool) {
	var input *dto.VerifyOTP
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.DefaultResponse{Message: err.Error()})
		return &dto.OTPValidationCredentials{}, false
	}
	cred := &dto.OTPValidationCredentials{}
	IPAddress, _ := helper.GetIPAddress(c)
	cred.Email, cred.OTP, cred.IPAddress = input.Email, input.OTP, IPAddress
	if !providers.EmailValidationLimiter.AllowRequest(cred.Email, cred.IPAddress) {
		c.JSON(http.StatusTooManyRequests, dto.DefaultResponse{Message: "Too Many OTP Validation Attempts, Retry Later"})
		return &dto.OTPValidationCredentials{}, false
	}
	return cred, true
}

func AccessTokenAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := guard.GetBearerToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.DefaultResponse{Message: err.Error()})
			return
		}
		err = providers.JWTMaker.VerifyToken(token, "access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.DefaultResponse{Message: err.Error()})
			return
		}
		c.Next()
	}
}

func RefreshTokenAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := guard.GetBearerToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.DefaultResponse{Message: err.Error()})
			return
		}
		err = providers.JWTMaker.VerifyToken(token, "refresh_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.DefaultResponse{Message: err.Error()})
			return
		}
		c.Next()
	}
}
