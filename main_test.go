package main

import (
	"log"
	"os"
	"testing"

	"github.com/conormkelly/fiber-demo/database"
	"github.com/conormkelly/fiber-demo/models"
)

var app App

// Create an in-memory SQLite DB for testing purposes.
func TestMain(m *testing.M) {
	db := database.Database{}
	modelsToMigrate := []interface{}{&models.User{}}
	err := db.Connect(&database.Options{UseInMemoryDatabase: true, ModelsToMigrate: modelsToMigrate})
	if err != nil {
		log.Fatal("Database failed to connect: " + err.Error())
	}

	app.Initialize(&db)

	code := m.Run()
	os.Exit(code)
}
