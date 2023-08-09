package models

import (
	"github.com/ignitedotdev/auth-ms/internal/api/common/entities/models"
	"gorm.io/gorm"
)

type Developer struct {
	gorm.Model
	models.BUser
}
