package routes

import (
	mid "YenExpress/service/middlewares"
	con "YenExpress/service/patient/controllers"
	p_mid "YenExpress/service/patient/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine) {

	group := router.Group("/patient/auth")

	group.POST("/register/", mid.APIKeyAuthorization(), con.Register)
	group.POST("/login/", mid.APIKeyAuthorization(), con.Login)
	group.POST("/email/send-otp", mid.APIKeyAuthorization(), con.SendNewMailOTP)
	group.POST("/email/verify", mid.APIKeyAuthorization(), con.ConfirmNewMail)
	group.DELETE("/logout/", mid.APIKeyAuthorization(), p_mid.AccessTokenAuthorization(), con.Logout)
	group.GET("/refresh/", mid.APIKeyAuthorization(), p_mid.RefreshTokenAuthorization(), con.Refresh)
}
