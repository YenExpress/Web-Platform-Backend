package models

import (
	"YenExpress/config"
	"encoding/json"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateDTO struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Role      string `json:"role" validate:"required,Enum=admin_superadmin"`
}

type Admin struct {
	gorm.Model
	FirstName string
	LastName  string
	Role      string `gorm:"default:Admin"`
	Email     string
	Password  string
	Sex       string `gorm:"default:Male"`
	Photo     string `gorm:"not null:false"`
}

func (user *Admin) SaveNew() (*Admin, error) {
	err := config.DB.Create(&user).Error
	if err != nil {
		return &Admin{}, err
	}
	return user, nil
}

func (user *Admin) BeforeCreate(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	return nil
}

func (user *Admin) ValidatePwd(provided_password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(provided_password))
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
