package entities

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Base user struct for defining user attributes in database
type BUser struct {
	gorm.Model
	FirstName      string `gorm:"not null:false"`
	LastName       string `gorm:"not null:false"`
	Email          string `gorm:"unique"`
	HashedPassword string
}

//  Hash User Password before saving to database
func (user *BUser) BeforeCreate(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.HashedPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.HashedPassword = string(passwordHash)
	return nil
}

// Compare password with hash value saved in database
func (user *BUser) ValidatePwd(strPwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(strPwd))
}

// Developer model builds on base user implemeting same methods
type Patient struct {
	gorm.Model
	BUser
}
