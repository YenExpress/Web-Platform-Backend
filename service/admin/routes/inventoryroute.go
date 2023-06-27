package routes

import (
	con "YenExpress/service/admin/controllers"

	"github.com/gin-gonic/gin"
)

func ProductRoute(router *gin.Engine) {

	group := router.Group("/admin/product")

	group.POST("/list-drug", con.ListDrug)
	group.POST("/add-category", con.AddDrugCategory)

}
