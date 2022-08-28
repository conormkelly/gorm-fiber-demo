package main

import (
	"os"
	"testing"

	"github.com/conormkelly/fiber-demo/database"
)

var app App

// Create an in-memory SQLite DB for testing purposes.
func TestMain(m *testing.M) {
	db := database.Database{}
	db.Connect(&database.Options{UseInMemoryDatabase: true, PerformMigration: true})

	app.Initialize(&db)

	code := m.Run()
	os.Exit(code)
}
