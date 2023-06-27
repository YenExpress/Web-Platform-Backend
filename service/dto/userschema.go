package dto

import (
	"YenExpress/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"

	"gorm.io/gorm"
)

type WaitList struct {
	gorm.Model
	Email string `gorm:"unique"`
}

func (list *WaitList) SaveNew() error {
	return config.DB.Create(&list).Error
}

type Patient struct {
	gorm.Model
	// ID          string `gorm:"primarykey"`
	FirstName   string
	LastName    string
	UserName    string
	Email       string `gorm:"unique"`
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
	// user.ID = uuid.New().String()
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

type Admin struct {
	gorm.Model
	// ID        string `gorm:"primarykey"`
	FirstName string
	LastName  string
	Role      string `gorm:"default:Admin"`
	Email     string `gorm:"unique"`
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
	// user.ID = uuid.New().String()
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
