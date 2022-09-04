package database

import (
	"testing"

	"github.com/conormkelly/fiber-demo/models"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryDB(t *testing.T) {
	db := Database{}
	modelsToMigrate := []interface{}{&models.User{}}
	connectionString := "file:db_test?mode=memory&cache=shared"

	err := db.Connect(&Options{
		DatabaseType:     SQLite,
		ConnectionString: &connectionString,
		ModelsToMigrate:  modelsToMigrate,
	})
	assert.Nil(t, err, "Should have connected successfuly but it returned an unexpected error")
}

func TestMissingConnectionString(t *testing.T) {
	db := Database{}
	err := db.Connect(&Options{
		DatabaseType:     SQLite,
		ConnectionString: nil,
	})
	assert.NotNil(t, err, "Expected an error but got none")
	assert.Equal(t, "no ConnectionString was provided", err.Error())
}

func TestInvalidConnectionString(t *testing.T) {
	db := Database{}
	connectionString := "/;"
	err := db.Connect(&Options{
		DatabaseType:     SQLite,
		ConnectionString: &connectionString,
	})
	assert.NotNil(t, err, "Expected an error but got none")
}

func TestUndefinedDatabaseType(t *testing.T) {
	db := Database{}
	connectionString := "file:TestUndefinedDatabaseType?mode=memory&cache=shared"
	err := db.Connect(&Options{
		ConnectionString: &connectionString,
	})
	assert.NotNil(t, err, "Expected an error but got none")
	assert.Equal(t, "no DatabaseType was specified", err.Error())
}
