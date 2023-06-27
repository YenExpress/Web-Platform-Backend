package dto

import (
	"time"

	"gorm.io/gorm"
)

type DrugOrder struct {
	gorm.Model
	OrderDate        time.Time
	ItemsDescription string
	AmountPaid       string
	Status           string
	DeliveryAddress  string
	CustomerPhone    string
	CustomerID       *uint    `gorm:"not null:false"`
	Customer         *Patient `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,ONDELETE:SET NULL;"`
}
