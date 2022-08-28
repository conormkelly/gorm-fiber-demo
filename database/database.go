package database

import (
	"log"

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
	ModelsToMigrate     []interface{}
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

	if len(options.ModelsToMigrate) > 0 {
		log.Println("Running migrations.")
		// Variadic function
		db.AutoMigrate(options.ModelsToMigrate...)
	}

	dbInstance.Conn = db

	return nil
}
