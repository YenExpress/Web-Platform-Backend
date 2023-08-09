package database

import (
	"fmt"

	"github.com/ignitedotdev/auth-ms/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect to database and create tables with defined models
func ConnectDB(models ...interface{}) {
	var err error

	database, err := gorm.Open(postgres.Open(config.Config.PostgresURI), &gorm.Config{})

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

// Delete all tables in database with defined models
func CleanDB(models ...interface{}) {
	migrator := DB.Migrator()
	for _, model := range models {
		migrator.DropTable(model)
		fmt.Println("Deleted All Records in ", model)
	}
}
