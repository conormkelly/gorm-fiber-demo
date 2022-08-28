package database

import (
	"log"

	"github.com/conormkelly/fiber-demo/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Conn *gorm.DB
}

type Options struct {
	UseInMemoryDatabase bool
	SQLitePath          *string
	PerformMigration    bool
}

func (dbInstance *Database) Connect(options *Options) error {
	var db *gorm.DB
	var err error

	if options.UseInMemoryDatabase {
		db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	} else if options.SQLitePath != nil {
		db, err = gorm.Open(sqlite.Open(*options.SQLitePath), &gorm.Config{})
	} else {
		log.Fatalln("Invalid DB config provided.")
	}

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err.Error())
	}

	log.Println("Connected to the database.")
	db.Logger = logger.Default.LogMode(logger.Info)

	if options.PerformMigration {
		log.Println("Running migrations...")
		// AutoMigrate is a variadic function,
		// so if you have multiple models, you can pass them all
		db.AutoMigrate(&models.User{})
	}

	dbInstance.Conn = db

	return nil
}
