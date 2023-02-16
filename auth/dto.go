package auth

import (
	"YenExpress/config"
	"YenExpress/guard"

	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

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
	LastLogin   mysql.NullTime
	IsActive    bool `gorm:"default:false"`
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

func (user *Patient) GetAccessToken() string {
	return guard.PatientJWTMaker.CreatePatientToken(time.Now().Add(time.Hour*24*3),
		user.ID, "access_token")
}

func (user *Patient) GetRefreshToken() string {
	return guard.PatientJWTMaker.CreatePatientToken(time.Now().Add(time.Hour*24*30),
		user.ID, "refresh_token")
}

func (user *Patient) GetIDToken() string {
	return guard.PatientJWTMaker.CreatePatientIdToken(user.ID,
		user.LastName, user.FirstName, user.Email, user.Sex)
}
