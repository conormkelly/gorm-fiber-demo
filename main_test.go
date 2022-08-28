package main

import (
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
	db.Connect(&database.Options{UseInMemoryDatabase: true, ModelsToMigrate: modelsToMigrate})

	app.Initialize(&db)

	code := m.Run()
	os.Exit(code)
}
