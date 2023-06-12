package dto

import (
	"YenExpress/helper"

	"github.com/gin-gonic/gin"
)

type LoginCredentials struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	IPAddress string
}

type ConfirmEmail struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTP struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=7"`
}

type OTPValidationCredentials struct {
	IPAddress string
	Email     string
	OTP       string
}

func (cred *OTPValidationCredentials) LoadFromParams(c *gin.Context) error {
	ipAddr, err := helper.GetIPAddress(c)
	if err != nil {
		return err
	}
	cred.IPAddress = ipAddr
	return nil

}
