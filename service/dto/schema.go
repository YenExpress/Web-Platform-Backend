package dto

import (
	"YenExpress/config"
	"time"

	"gorm.io/gorm"
)

type WaitList struct {
	gorm.Model
	Email string `gorm:"unique"`
}

func (list *WaitList) SaveNew() error {
	return config.DB.Create(&list).Error
}

type Drug struct {
	gorm.Model
	ID           string `gorm:"unique"`
	BrandName    string
	Category     string
	Price        string
	Description  string
	Availability string
	Photo        string `gorm:"not null:false"`
}

type DrugOrder struct {
	gorm.Model
	ID               string `gorm:"unique"`
	OrderDate        time.Time
	ItemsDescription string
	AmountPaid       string
	Status           string
	DeliveryAddress  string
	CustomerPhone    string
	CustomerID       string
}
