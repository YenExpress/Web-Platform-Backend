package interactors

import (
	"errors"

	"github.com/ignitedotdev/auth-ms/internal/api/database/entities"

	"github.com/ignitedotdev/auth-ms/internal/api/common/exceptions"
	repo "github.com/ignitedotdev/auth-ms/internal/api/shared/repositories"

	"gorm.io/gorm"
)

// Use case object to handle all developer related authentication processes
// Service hinged upon repository serving as data access layer
type DeveloperAuthService struct {
	repository repo.IUserRepository[entities.Developer]
}

// constructor function to create an instance of developer auth use case object
func NewDeveloperAuthService(userRepository repo.IUserRepository[entities.Developer]) *DeveloperAuthService {
	return &DeveloperAuthService{repository: userRepository}
}

// register user on platform with required information provided
func (usecase *DeveloperAuthService) NativeSignUp(firstName, lastName, email, password string) error {

	newDeveloper := &entities.Developer{
		BUser: entities.BUser{
			FirstName: firstName, LastName: lastName,
			Email: email, HashedPassword: password,
		},
	}

	if err := usecase.repository.SaveNew(newDeveloper); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return exceptions.UserExists
		}
		return err
	}
	return nil

}

// sign use in with credentails `email` and `password`
func (usecase *DeveloperAuthService) NativeLogin(email, password string) error {

	developer, err := usecase.repository.GetByEmail(email)
	if err != nil {
		return err
	} else if developer == nil {
		return exceptions.UserDoesNotExist
	}

	if err = developer.ValidatePwd(password); err != nil {
		return exceptions.InvalidPassword
	}
	return nil

}
