// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import ("time")

type Drug struct {
	ID           string  `json:"id"`
	BrandName    string  `json:"brandName"`
	Category     string  `json:"category"`
	Price        string  `json:"price"`
	Description  *string `json:"description,omitempty"`
	Availability *string `json:"availability,omitempty"`
	Photo *string `json:"photo,omitempty"`
}

type DrugOrder struct {
	ID               string  `json:"id"`
	OrderDate        time.Time  `json:"orderDate"`
	ItemsDescription *string `json:"itemsDescription,omitempty"`
	AmountPaid       string  `json:"amountPaid"`
	Status           string  `json:"status"`
	DeliveryAddress  string  `json:"deliveryAddress"`
	CustomerPhone    string  `json:"customerPhone"`
	CustomerID       string  `json:"customerID"`
}
