package models

import (
	"YenExpress/config"
	"encoding/json"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

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

type RegistrationDTO struct {
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	Sex         string `json:"sex" validate:"required,Enum=male_female"`
	Location    string `json:"location,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty" binding:"e164"`
}

type SessionData struct {
	SessionID string    `json:"sessionID"`
	IPAddress string    `json:"ipAddr,omitempty"`
	Email     string    `json:"Email"`
	LoggedIn  time.Time `json:"loggedIn"`
	UserID    uint      `json:"userID"`
}

type SessionStore struct {
	Session map[string]SessionData
}

func (s *SessionStore) Encode() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SessionStore) Decode(data []byte) error {
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return nil
}
