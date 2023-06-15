package middlewares

import (
	"YenExpress/helper"
	pro "YenExpress/service/admin/providers"
	"YenExpress/service/dto"
	"YenExpress/service/guard"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RateLimitLogin(c *gin.Context) (*dto.LoginCredentials, bool) {
	var input *dto.LoginCredentials
	IPAddress, _ := helper.GetIPAddress(c)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.DefaultResponse{Message: err.Error()})
		return &dto.LoginCredentials{}, false
	}
	if !pro.LoginLimiter.AllowRequest(input.Email, IPAddress) {
		c.JSON(http.StatusTooManyRequests, dto.DefaultResponse{Message: "Too Many Failed Login Attempts, Retry Later"})
		return &dto.LoginCredentials{}, false
	}
	input.IPAddress = IPAddress
	return input, true
}

func RateLimitOTPGeneration(c *gin.Context) (*dto.OTPValidationCredentials, bool) {
	var input *dto.ConfirmEmail
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.DefaultResponse{Message: err.Error()})
		return &dto.OTPValidationCredentials{}, false
	}
	cred := &dto.OTPValidationCredentials{}
	IPAddress, _ := helper.GetIPAddress(c)
	cred.Email, cred.IPAddress = input.Email, IPAddress
	if !pro.CreateOTPLimiter.AllowRequest(cred.Email, cred.IPAddress) {
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
	if !pro.EmailValidationLimiter.AllowRequest(cred.Email, cred.IPAddress) {
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
		err = pro.JWTMaker.VerifyToken(token, "access_token")
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
		err = pro.JWTMaker.VerifyToken(token, "refresh_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.DefaultResponse{Message: err.Error()})
			return
		}
		c.Next()
	}
}
