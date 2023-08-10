package interactors

import (
	"errors"

	EN "github.com/ignitedotdev/auth-ms/internal/api/database/entities"

	EXC "github.com/ignitedotdev/auth-ms/internal/api/common/exceptions"
	SI "github.com/ignitedotdev/auth-ms/internal/api/shared/interactor"
	INT "github.com/ignitedotdev/auth-ms/internal/api/shared/interfaces"

	"gorm.io/gorm"
)

// Use case object to handle all developer related authentication processes
// Service hinged upon repository serving as data access layer
type DeveloperAuthService struct {
	SI.AuthService
	repository INT.IRepository[EN.Developer]
}

// constructor function to create an instance of developer auth use case object
func NewDeveloperAuthService(userRepository INT.IRepository[EN.Developer]) *DeveloperAuthService {
	return &DeveloperAuthService{repository: userRepository}
}

// register user on platform with required information provided
func (usecase *DeveloperAuthService) NativeSignUp(firstName, lastName, email, password string) error {

	newDeveloper := &EN.Developer{
		BUser: EN.BUser{
			FirstName: firstName, LastName: lastName,
			Email: email, HashedPassword: password,
		},
	}

	if err := usecase.repository.SaveNew(newDeveloper); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return EXC.UserExists
		}
		return err
	}
	return nil

}

// // sign use in with credentails `email` and `password`
// func (usecase *DeveloperAuthService) NativeLogin(email, password string) error {

// 	developer, err := usecase.repository.GetByEmail(email)
// 	if err != nil {
// 		return err
// 	} else if developer == nil {
// 		return EXC.UserDoesNotExist
// 	}

// 	if err = developer.ValidatePwd(password); err != nil {
// 		return EXC.InvalidPassword
// 	}
// 	return nil

// }
