package controllers

import (
	"YenExpress/config"
	"YenExpress/service/admin/models"
	"YenExpress/service/dto"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Save New Drug entry to database for drug listings
func ListDrug(c *gin.Context) {

	func() {

		var input *models.CreateDrugDTO

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, dto.DefaultResponse{Message: err.Error()})
			return
		}

		var drug *dto.Drug
		err := config.DB.Where("brandName = ?", input.BrandName).First(&drug).Error
		if err == nil {
			c.JSON(http.StatusConflict, dto.DefaultResponse{Message: "Drug Brand Already Listed"})
			return
		}

		drug = &dto.Drug{
			BrandName: input.BrandName, CategoryID: &input.CategoryID,
			Photo: input.Photo, Description: input.Description,
			Availability: input.Availability, Price: input.Price,
		}

		err = drug.SaveNew()

		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.DefaultResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, dto.DefaultResponse{Message: fmt.Sprintf("%s Drug Listed", input.BrandName)})
		return

	}()
}

// Create new drug category for drug listings
func AddDrugCategory(c *gin.Context) {

	func() {

		var input *models.CreateDrugCategoryDTO

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, dto.DefaultResponse{Message: err.Error()})
			return
		}

		var category *dto.DrugCategory
		err := config.DB.Where("Name = ?", input.Name).First(&category).Error
		if err == nil {
			c.JSON(http.StatusConflict, dto.DefaultResponse{Message: fmt.Sprintf("Category Named %s Already Listed", input.Name)})
			return
		}

		category = &dto.DrugCategory{
			Name: input.Name, ParentCategoryID: input.ParentCategoryID,
			Description: input.Description,
		}

		err = category.SaveNew()

		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.DefaultResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, dto.DefaultResponse{Message: fmt.Sprintf("Another Category %s Created ", input.Name)})
		return

	}()
}
