package auth

import (
	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine) {

	group := router.Group("/patient/auth")

	group.POST("/create/", Register)
	group.POST("/login/", Login)
	group.GET("/sendotp/:email/*process", GetOneTimePass)
	group.GET("/validateotp/:email/:otp/*process", ValidateOneTimePass)
}
