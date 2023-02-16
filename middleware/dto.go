package middleware

type DefaultResponse struct {
	Message string `json:"message,omitempty"`
}

type LoginCredentials struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	IPAddress string
}
