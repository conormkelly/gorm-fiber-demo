package database

import (
	"os"
	"testing"
)

// Quick sanity check for the package
func TestMain(m *testing.M) {
	db := Database{}
	db.Connect(&Options{UseInMemoryDatabase: true, PerformMigration: true})

	code := m.Run()
	os.Exit(code)
}
