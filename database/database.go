package database

import (
	"errors"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Conn *gorm.DB
}

type Options struct {
	ConnectionString *string
	ModelsToMigrate  []interface{}
}

func GetConnection(options *Options) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	if options.ConnectionString == nil {
		return nil, errors.New("no ConnectionString was provided")
	}

	db, err = gorm.Open(mysql.Open(*options.ConnectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Connected to the database.")
	db.Logger = logger.Default.LogMode(logger.Info)

	if len(options.ModelsToMigrate) > 0 {
		log.Println("Running migrations.")
		// Variadic function
		err = db.AutoMigrate(options.ModelsToMigrate...)
	}

	return db, err
}
