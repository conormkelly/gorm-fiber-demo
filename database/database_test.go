package database

import (
	"testing"

	"github.com/conormkelly/fiber-demo/models"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryDB(t *testing.T) {
	db := Database{}
	modelsToMigrate := []interface{}{&models.User{}}
	err := db.Connect(&Options{UseInMemoryDatabase: true, ModelsToMigrate: modelsToMigrate})
	assert.Nil(t, err, "Should have connected successfuly but it returned an unexpected error")
}

func TestMissingPath(t *testing.T) {
	db := Database{}
	err := db.Connect(&Options{UseInMemoryDatabase: false})
	assert.NotNil(t, err, "Expected an error but got none")
	assert.Equal(t, "invalid DB config provided", err.Error())
}

func TestInvalidPath(t *testing.T) {
	db := Database{}
	connectionString := "?"
	err := db.Connect(&Options{UseInMemoryDatabase: false, SQLitePath: &connectionString})
	assert.NotNil(t, err, "Expected an error but got none")
	assert.Equal(t, "unable to open database file: The filename, directory name, or volume label syntax is incorrect.", err.Error())
}
