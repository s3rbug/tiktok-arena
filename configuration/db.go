package configuration

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"tiktok-arena/models"
)

var DB *gorm.DB

func ConnectDB(config *EnvConfigModel) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.DBHost,
		config.DBUserName,
		config.DBUserPassword,
		config.DBName,
		config.DBPort,
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database!\n", err.Error())
	}
	//	Extension for postgresql uuid support
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Tournament{},
		&models.Tiktok{},
	)
	if err != nil {
		log.Fatal("Migration Failed:\n", err.Error())
	}
	log.Println("Successfully connected to the database")
}
