package interfaces

// interface describing common repository behaviours for different user models
type IUserRepository[userSchema interface{}] interface {
	SaveNew(user *userSchema) error
	GetByEmail(email string) (*userSchema, error)
	GetByID(ID interface{}) (*userSchema, error)
	Exists(ID interface{}) bool
}
