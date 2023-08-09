package repositories

import (
	"github.com/ignitedotdev/auth-ms/internal/database"
)

// interface describing common repository behaviours for different user models
type IUserRepository[userSchema interface{}] interface {
	SaveNew(user *userSchema) error
	GetByEmail(email string) (*userSchema, error)
	GetByID(ID interface{}) (*userSchema, error)
	Exists(ID interface{}) bool
}

// base user repository exposing common methods for child user repositories to override/inherit
type BUserRepository[userSchema interface{}] struct {
}

func (repo *BUserRepository[userSchema]) SaveNew(user *userSchema) error {
	err := database.DB.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// get user record from database based on email
func (repo *BUserRepository[userSchema]) GetByEmail(email string) (*userSchema, error) {
	var schema *userSchema
	err := database.DB.Where("email = ?", email).First(&schema).Error
	if err == nil {
		return schema, nil
	}
	return nil, err
}

// get user record from database based on id
func (repo *BUserRepository[userSchema]) GetByID(ID interface{}) (*userSchema, error) {
	var schema *userSchema
	err := database.DB.Where("ID = ?", ID).First(&schema).Error
	if err == nil {
		return schema, nil
	}
	return nil, err
}

// check if user with specified ID exists in record base
func (repo *BUserRepository[userSchema]) Exists(ID interface{}) bool {
	var schema *userSchema
	err := database.DB.Where("ID = ?", ID).First(&schema).Error
	if err != nil {
		return false
	}
	return true
}
