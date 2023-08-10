package repository

import (
	database "github.com/ignitedotdev/auth-ms/internal/api/database/connectors"
)

// base user repository exposing common methods for child user repositories to override/inherit
type BRepository[userSchema interface{}] struct {
}

func (repo *BRepository[userSchema]) SaveNew(user *userSchema) error {
	err := database.DB.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// get user record from database based on email
func (repo *BRepository[userSchema]) GetByEmail(email string) (*userSchema, error) {
	var schema *userSchema
	err := database.DB.Where("email = ?", email).First(&schema).Error
	if err == nil {
		return schema, nil
	}
	return nil, err
}

// get user record from database based on id
func (repo *BRepository[userSchema]) GetByID(ID interface{}) (*userSchema, error) {
	var schema *userSchema
	err := database.DB.Where("ID = ?", ID).First(&schema).Error
	if err == nil {
		return schema, nil
	}
	return nil, err
}

// check if user with specified ID exists in record base
func (repo *BRepository[userSchema]) Exists(ID interface{}) bool {
	var schema *userSchema
	err := database.DB.Where("ID = ?", ID).First(&schema).Error
	if err != nil {
		return false
	}
	return true
}
