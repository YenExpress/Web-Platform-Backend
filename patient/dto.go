package patient

// import (

// 	"time"

// 	"golang.org/x/crypto/bcrypt"
// 	"gorm.io/datatypes"
// 	"gorm.io/gorm"
// )

// type Patient struct {
// 	gorm.Model
// 	ID           uint `gorm:"primary_key"`
// 	FirstName    string
// 	LastName     string
// 	LastLogin    time.Time `gorm:"not null:false"`
// 	LastModified time.Time `gorm:"not null:false"`
// 	IsActive     bool      `gorm:"not null:false"`
// 	Email        string
// 	Password     string
// 	Created      time.Time
// 	DateOfBirth  datatypes.Date `gorm:"not null:false"`
// 	Sex          string         `gorm:"not null:true;default:male"`
// 	Photo        string         `gorm:"not null:false"`
// 	Location     string         `gorm:"not null:false"`
// 	PhoneNumber  string         `gorm:"not null:false"`
// }

// type ExportDTO struct {
// 	FirstName string `json:"firstName,omitempty" validate:"required"`
// 	LastName  string `json:"lastName,omitempty" validate:"required"`
// 	Email     string `json:"email,omitempty" validate:"required"`
// 	Sex       string `json:"sex,omitempty" validate:"required"`
// }

// type EditBioDTO struct {
// 	FirstName string `json:"firstName,omitempty" validate:"required"`
// 	LastName  string `json:"lastName,omitempty" validate:"required"`
// 	Email     string `json:"email,omitempty" validate:"required"`
// 	Sex       string `json:"sex,omitempty" validate:"required"`
// 	Photo     string `json:"photo,omitempty" validate:"required"`
// }
