package dto

type LoginCredentials struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	IPAddress string
}

type ConfirmEmail struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTP struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=7"`
}

type OTPValidationCredentials struct {
	IPAddress string
	Email     string
	OTP       string
}
