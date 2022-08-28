package database

import (
	"os"
	"testing"

	"github.com/conormkelly/fiber-demo/models"
)

// Quick sanity check for the package
func TestMain(m *testing.M) {
	db := Database{}
	modelsToMigrate := []interface{}{&models.User{}}
	db.Connect(&Options{UseInMemoryDatabase: true, ModelsToMigrate: modelsToMigrate})

	code := m.Run()
	os.Exit(code)
}
