package dto

import (
	"YenExpress/config"

	"gorm.io/gorm"
)

type WaitList struct {
	gorm.Model
	Email string `gorm:"unique"`
}

func (list *WaitList) SaveNew() error {
	return config.DB.Create(&list).Error
}
