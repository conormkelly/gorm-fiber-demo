package database

import (
	"errors"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Conn *gorm.DB
}

type Options struct {
	UseInMemoryDatabase  bool
	InMemoryDatabaseName string
	SQLitePath           *string
	ModelsToMigrate      []interface{}
}

func (dbInstance *Database) Connect(options *Options) error {
	var db *gorm.DB
	var err error

	if options.UseInMemoryDatabase {
		connectionString := "file::memory:?cache=shared"
		if options.InMemoryDatabaseName != "" {
			connectionString = "file:" + options.InMemoryDatabaseName + "?mode=memory&cache=shared"
		}
		db, err = gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
	} else if options.SQLitePath != nil {
		db, err = gorm.Open(sqlite.Open(*options.SQLitePath), &gorm.Config{})
	} else {
		return errors.New("invalid DB config provided")
	}

	if err != nil {
		return err
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
