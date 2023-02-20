package auth

import (
	pAuth "YenExpress/auth/patientauth"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine) {

	group := router.Group("/auth")

	group.POST("/create/patient", pAuth.RegisterPatient)
	group.POST("/login/patient", pAuth.LoginPatient)
}
