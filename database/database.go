package database

import (
	"log"

	"github.com/conormkelly/fiber-demo/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database SQLDatabase

type SQLDatabase struct {
	DB *gorm.DB
}

func (dbInstance SQLDatabase) Connect() error {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err.Error())
	}

	log.Println("Connected to the database.")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running migrations...")
	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})

	Database.DB = db

	return nil
}
