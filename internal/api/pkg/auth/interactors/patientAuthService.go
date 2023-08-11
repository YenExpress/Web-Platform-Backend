package interactors

import (
	"errors"

	EN "yenexpress/internal/api/database/entities"

	EXC "yenexpress/internal/api/pkg/auth/exceptions"
	INT "yenexpress/internal/api/pkg/shared/interfaces"

	"gorm.io/gorm"
)

// Use case object to handle all patient related authentication processes
// Service hinged upon repository serving as data access layer
type PatientAuthService struct {
	AuthService
	repository INT.IUserRepository[EN.Patient]
}

// constructor function to create an instance of developer auth use case object
func NewPatientAuthService(userRepository INT.IUserRepository[EN.Patient]) *PatientAuthService {
	return &PatientAuthService{repository: userRepository}
}

// register user on platform with required information provided
func (usecase *PatientAuthService) NativeSignUp(firstName, lastName, email, password string) error {

	newDeveloper := &EN.Patient{
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
