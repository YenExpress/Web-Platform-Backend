package dto

import (
	"YenExpress/config"
	"gorm.io/gorm"
)

type DrugCategory struct {
	gorm.Model
	Name             string
	Description      string `gorm:"not null:false"`
	ParentCategoryID *uint
	ParentCategory   *DrugCategory `gorm:"foreignKey:ParentCategoryID;constraint:OnUpdate:CASCADE,ONDELETE:SET NULL;"`
}

func (item *DrugCategory) SaveNew() error {
	err := config.DB.Create(&item).Error
	return err
}

type Drug struct {
	gorm.Model
	BrandName    string
	CategoryID   *uint
	Category     *DrugCategory `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,ONDELETE:SET NULL;"`
	Price        string
	Description  string
	Availability string `gorm:"not null:false"`
	Photo        string `gorm:"not null:false"`
}

// func (item *Drug) BeforeCreate(*gorm.DB) error {
// 	if id == "" {
// 		return errors.New("Invalid UUID")
// 	}
// 	item.ID = id
// 	return nil
// }

func (item *Drug) SaveNew() error {
	err := config.DB.Create(&item).Error
	return err
}
