package exceptions

import "errors"

var (
	UserExists       = errors.New("Account Already Exists")
	InvalidEmail     = errors.New("Email Provided Incorrect")
	InvalidPassword  = errors.New("Password Incorrect")
	UserDoesNotExist = errors.New("Account Does Not Exist")
	ErrInvalidToken  = errors.New("Invalid or unrecognized Token")
	ErrExpiredToken  = errors.New("Token Expired")
	ErrInvalidAPIKey = errors.New("Invalid API Key")
)
