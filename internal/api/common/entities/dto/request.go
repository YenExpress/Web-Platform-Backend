package dto

type LoginCredentials struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	IPAddress string
}
