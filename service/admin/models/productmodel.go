package models

type CreateDrugDTO struct {
	BrandName    string `json:"brandName" binding:"required"`
	Price        string `json:"price" binding:"required"`
	Description  string `json:"description" binding:"required"`
	CategoryID   uint   `json:"categoryID" binding:"required"`
	Photo        string `json:"photo"`
	Availability string `json:"availability"`
}

type CreateDrugCategoryDTO struct {
	Name             string `json:"name" binding:"required"`
	Description      string `json:"description" binding:"required"`
	ParentCategoryID *uint  `json:"parentcategoryID"`
}
