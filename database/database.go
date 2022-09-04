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
	ConnectionString *string
	DatabaseType     DatabaseType
	ModelsToMigrate  []interface{}
}

// SQLite, MySQL etc
type DatabaseType int32

const (
	Undefined DatabaseType = iota
	SQLite
	// MySQL
)

func (dbInstance *Database) Connect(options *Options) error {
	var db *gorm.DB
	var err error

	if options.ConnectionString == nil {
		return errors.New("no ConnectionString was provided")
	}

	switch options.DatabaseType {
	case Undefined:
		return errors.New("no DatabaseType was specified")
	case SQLite:
		db, err = gorm.Open(sqlite.Open(*options.ConnectionString), &gorm.Config{})
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
