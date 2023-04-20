package auth

import (
	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine) {

	group := router.Group("/patient/auth")

	group.POST("/register/", AuthorizeWithAPIKey(), Register)
	group.POST("/login/", AuthorizeWithAPIKey(), Login)
	group.POST("/confirm-email/:process/", AuthorizeWithAPIKey(), GenerateOTPforAuth)
	group.POST("/verify-otp/:process/", AuthorizeWithAPIKey(), ValidateOneTimePass)
}
