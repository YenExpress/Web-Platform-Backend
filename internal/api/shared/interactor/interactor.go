package interactor

import (
	EXC "github.com/ignitedotdev/auth-ms/internal/api/common/exceptions"
	EN "github.com/ignitedotdev/auth-ms/internal/api/database/entities"
	INT "github.com/ignitedotdev/auth-ms/internal/api/shared/interfaces"
)

// Use case object to handle all developer related authentication processes
// Service hinged upon repository serving as data access layer
type AuthService struct {
	repository INT.IRepository[EN.BUser]
}

// // constructor function to create an instance of developer auth use case object
// func NewDeveloperAuthService(userRepository INT.IRepository[EN.Developer]) *DeveloperAuthService {
// 	return &DeveloperAuthService{repository: userRepository}
// }

// sign use in with credentails `email` and `password`
func (usecase *AuthService) NativeLogin(email, password string) error {

	developer, err := usecase.repository.GetByEmail(email)
	if err != nil {
		return err
	} else if developer == nil {
		return EXC.UserDoesNotExist
	}

	if err = developer.ValidatePwd(password); err != nil {
		return EXC.InvalidPassword
	}
	return nil

}
