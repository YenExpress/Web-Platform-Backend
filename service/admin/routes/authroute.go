package routes

import (
	con "YenExpress/service/admin/controllers"
	admin_mid "YenExpress/service/admin/middlewares"
	mid "YenExpress/service/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine) {

	group := router.Group("/admin/auth")

	group.POST("/create-account", mid.APIKeyAuthorization(), con.CreateAdmin)
	group.POST("/login", mid.APIKeyAuthorization(), con.Login)
	group.POST("/login/send-otp", mid.APIKeyAuthorization(), con.GenerateAuthOTP)
	group.POST("/login/verify", mid.APIKeyAuthorization(), con.ValidateLoginOTP)
	group.DELETE("/logout", mid.APIKeyAuthorization(), admin_mid.AccessTokenAuthorization(), con.Logout)
	group.GET("/refresh", mid.APIKeyAuthorization(), admin_mid.RefreshTokenAuthorization(), con.RefreshAuthToken)
}
