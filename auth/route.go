package auth

import (
	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine) {

	group := router.Group("/auth")

	group.POST("/create/patient", RegisterPatient)
	group.POST("/login/patient", LoginPatient)
}
