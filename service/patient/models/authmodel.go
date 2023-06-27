package models

import ()

type RegistrationDTO struct {
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	Sex         string `json:"sex" validate:"required,Enum=male_female"`
	Location    string `json:"location,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty" binding:"e164"`
}
