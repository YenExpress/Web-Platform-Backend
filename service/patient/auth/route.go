package auth

import (
	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine) {

	group := router.Group("/patient/auth")

	group.POST("/create/", AuthorizeWithAPIKey(), Register)
	group.POST("/login/", AuthorizeWithAPIKey(), Login)
	group.GET("/sendotp/:email/*process", AuthorizeWithAPIKey(), GetOneTimePass)
	group.GET("/validateotp/:email/:otp/*process", AuthorizeWithAPIKey(), ValidateOneTimePass)
}
