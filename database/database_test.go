package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingConnectionString(t *testing.T) {
	_, err := GetConnection(&Options{
		ConnectionString: nil,
	})
	assert.NotNil(t, err, "Expected an error but got none")
	assert.Equal(t, "no ConnectionString was provided", err.Error())
}

func TestInvalidConnectionString(t *testing.T) {
	connectionString := "/;"
	_, err := GetConnection(&Options{
		ConnectionString: &connectionString,
	})
	assert.NotNil(t, err, "Expected an error but got none")
}
