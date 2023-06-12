package main

import (
	"YenExpress/config"
	"YenExpress/docs"
	ad_model "YenExpress/service/admin/models"
	ad_route "YenExpress/service/admin/routes"
	p_model "YenExpress/service/patient/models"
	p_route "YenExpress/service/patient/routes"

	"YenExpress/helper"
	"fmt"
	"log"
	"net/http"
	"time"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {

	docs.SwaggerInfo.Title = "Yen Express APIs"
	docs.SwaggerInfo.Description = "This is the APIs server for the Yen Express Telemedicine Platform."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = config.ServerDomain
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	config.ConnectDB(&p_model.Patient{})
	config.ConnectDB(&ad_model.Admin{})
	helper.StartTaskMaster()

}

func main() {

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.DevOrigin, config.StagingOrigin, config.ProdOrigin},
		AllowMethods:     []string{},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour}))

	p_route.AuthRoute(router)
	ad_route.AuthRoute(router)

	router.GET("/", func(c *gin.Context) {
		ip, _ := helper.GetIPAddress(c)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Index Backend Service URL for YenExpress called by client with IP Address %v", ip),
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL(fmt.Sprintf("%v/swagger/doc.json", config.ServerDomain)),
		ginSwagger.DefaultModelsExpandDepth(1)))

	port := config.ServicePort
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
