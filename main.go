package main

import (
	"YenExpress/config"
	ad_model "YenExpress/service/admin/models"
	ad_route "YenExpress/service/admin/routes"
	"YenExpress/service/dto"
	p_model "YenExpress/service/patient/models"
	p_route "YenExpress/service/patient/routes"
	"YenExpress/service/searchAPI"

	"YenExpress/helper"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {

	config.ConnectDB(&p_model.Patient{})
	config.ConnectDB(&ad_model.Admin{})
	config.ConnectDB(&dto.Drug{})
	config.ConnectDB(&dto.DrugOrder{})
	config.ConnectDB(&dto.WaitList{})
	helper.StartTaskMaster()

}

func main() {

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.DevOrigin, config.StagingOrigin, config.ProdOrigin},
		AllowMethods:     []string{},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "x-api-key"},
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

	router.POST("/join-waitlist", func(c *gin.Context) {

		var input *dto.ConfirmEmail

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, dto.DefaultResponse{Message: err.Error()})
			return
		}

		var list *dto.WaitList
		err := config.DB.Where("email = ?", input.Email).First(&list).Error
		if err == nil {
			c.JSON(http.StatusOK, dto.DefaultResponse{Message: "Already on WaitList"})
			return
		}

		list = &dto.WaitList{Email: input.Email}

		err = list.SaveNew()

		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.DefaultResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, dto.DefaultResponse{Message: "WaitList Joined"})
		return

	})

	router.POST("/query", searchAPI.GraphqlHandler())
	router.GET("/graphql-playground", searchAPI.PlaygroundHandler())

	port := config.ServicePort
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
