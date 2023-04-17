package drugpurchase

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DefaultResponse struct {
	Message string `json:"message,omitempty"`
}

type DrugSalesPayment struct {
	gorm.Model
	OrderRef string
	Amount   string
	Provider string
	Status   string
}

type makeDrugSalesPaymentDTO struct {
	Amount   string `json:"amount" binding:"required"`
	Provider string `json:"provider" binding:"required"`
	Status   string `json:"status" binding:"required"`
}

type drugOrderDetails struct {
	ProductId   string     `json:"productId" binding:"required"`
	Description *time.Time `json:"description" binding:"required"`
	Amount      string     `json:"amount" binding:"required"`
}

type RetainDrugOrderDTO struct {
	OrderRef        string             `json:"orderRef" binding:"required"`
	Details         []drugOrderDetails `json:"details" binding:"required"`
	DeliveryAddress string             `json:"deliveryAddress" binding:"required"`
	CustomerPhone   string             `json:"customerPhone,omitempty" binding:"e164"`
	CustomerId      uint               `json:"customerId"`
}

type DrugOrder struct {
	gorm.Model
	OrderRef        string
	ProductId       string
	Description     string
	Amount          string
	DeliveryAddress string
	CustomerPhone   string
	CustomerId      uint `gorm:"not null:false"`
}

type DrugSales struct {
	gorm.Model
	OrderRef        string
	DateSold        datatypes.Date
	Status          string
	DateDelivered   datatypes.Date
	CustomerPhone   string
	DeliveryAddress string
	CustomerId      uint `gorm:"not null:false"`
}
