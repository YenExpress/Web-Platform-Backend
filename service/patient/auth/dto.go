package auth

import (
	"YenExpress/config"
	"encoding/json"
	"errors"
	"strings"

	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type confirmEmail struct {
	Email string `json:"email" binding:"required,email"`
}

type verifyOTP struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=7"`
}

type OTPValidationCredentials struct {
	ipAddress string
	process   string
	email     string
	otp       string
}

func (cred *OTPValidationCredentials) loadFromParams(c *gin.Context) error {
	ipAddr, err := config.GetIPAddress(c)
	if err != nil {
		return err
	}
	prcs := strings.Split(c.Param("process"), "/")[1]
	if prcs != "signin" && prcs != "signup" {
		return errors.New("process param value must be either signup or signin")
	}
	cred.ipAddress, cred.process = ipAddr, prcs
	return nil

}

type DefaultResponse struct {
	Message string `json:"message,omitempty"`
}

type LoginCredentials struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	IPAddress string
}

type LoginResponse struct {
	IDToken      string `json:"idToken"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
}

type CreatePatientDTO struct {
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	Sex         string `json:"sex" validate:"required,Enum=male_female"`
	Location    string `json:"location,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty" binding:"e164"`
}

type Patient struct {
	gorm.Model
	FirstName   string
	LastName    string
	UserName    string
	Email       string
	Password    string
	DateOfBirth datatypes.Date `gorm:"not null:false"`
	Sex         string         `gorm:"default:Male"`
	Photo       string         `gorm:"not null:false"`
	Location    string         `gorm:"not null:false"`
	PhoneNumber string         `gorm:"not null:false"`
}

func (user *Patient) SaveNew() (*Patient, error) {
	err := config.DB.Create(&user).Error
	if err != nil {
		return &Patient{}, err
	}
	return user, nil
}

func (user *Patient) BeforeCreate(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	return nil
}

func (user *Patient) ValidatePwd(provided_password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(provided_password))
}

type sessionData struct {
	SessionID string    `json:"sessionID"`
	IPAddress string    `json:"ipAddr,omitempty"`
	Email     string    `json:"Email"`
	LoggedIn  time.Time `json:"loggedIn"`
	UserID    uint      `json:"userID"`
}

type sessionStore struct {
	Session map[string]sessionData
}

func (s *sessionStore) encode() ([]byte, error) {
	return json.Marshal(s)
}

func (s *sessionStore) decode(data []byte) error {
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return nil
}
