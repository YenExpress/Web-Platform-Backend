package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(models ...interface{}) {
	var err error

	database, err := gorm.Open(postgres.Open(DatabaseURI), &gorm.Config{})

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully connected to the database")
	}
	for _, model := range models {
		database.AutoMigrate(model)
	}

	DB = database
}
