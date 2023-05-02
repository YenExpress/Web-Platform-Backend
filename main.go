package main

import (
	"YenExpress/config"
	"YenExpress/docs"
	patientAuth "YenExpress/service/patient/auth"
	"YenExpress/toolbox"
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

	config.ConnectDB(&patientAuth.Patient{})
	toolbox.StartTaskMaster()

}

func main() {

	router := gin.Default()

	patientAuth.AuthRoute(router)

	router.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{config.WebClientDomain},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PATCH", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin, Authorization"},
		ExposeHeaders:    []string{"*"},
		AllowAllOrigins:  true,
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == config.WebClientDomain
		// },
		MaxAge: 12 * time.Hour,
	}))
	router.GET("/", func(c *gin.Context) {
		ip, _ := config.GetIPAddress(c)
		fmt.Printf("client IP Address %v", ip)
		c.JSON(http.StatusOK, gin.H{
			"message": "Backend Service for YenExpress Telemedicine platform Healthy and Active.",
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
