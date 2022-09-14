package services

import (
	"fmt"
	"log"
	"testing"

	"github.com/conormkelly/fiber-demo/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type databaseTest struct {
	description string                   // Description of the test case
	action      func(*UserService) error // A func that relays the error from the service method
}

// Confirm that a specific error is returned by the service method when DB is unavailable
func TestDatabaseDown(t *testing.T) {
	testCases := []databaseTest{
		{
			description: "CreateUser with DB offline",
			action: func(svc *UserService) error {
				_, err := svc.CreateUser("Joe", "Bloggs")
				return err
			},
		},
		{
			description: "GetAllUsers with DB offline",
			action: func(svc *UserService) error {
				_, err := svc.GetAllUsers()
				return err
			},
		},
		{
			description: "GetUser with DB offline",
			action: func(svc *UserService) error {
				_, err := svc.GetUser(1)
				return err
			},
		},
		{
			description: "UpdateUser with DB offline",
			action: func(svc *UserService) error {
				firstName := "Joe"
				lastName := "Bloggs"
				_, err := svc.UpdateUser(1, &firstName, &lastName)
				return err
			},
		},
		{
			description: "DeleteUser with DB offline",
			action: func(svc *UserService) error {
				return svc.DeleteUser(1)
			},
		},
	}

	// Test setup / arrangement
	connectionString := "file:user_service_test?mode=memory&cache=shared"
	conn, err := gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
	db := &database.Database{Conn: conn}

	if err != nil {
		log.Fatal("Database failed to connect: " + err.Error())
	}

	// Close the database so that the calls to it in the service fail
	databaseConnection, err := db.Conn.DB()
	if err != nil {
		log.Fatal("Failed to get DB connection: " + err.Error())
	}
	databaseConnection.Close()

	userService := &UserService{DB: db}

	// Act
	expectedErrorMessage := "sql: database is closed"
	executeDbTests(userService, t, testCases, &expectedErrorMessage)
}

// Confirm that a specific error is returned by the service methods when table doesn't exist
func TestNonExistentTable(t *testing.T) {
	testCases := []databaseTest{
		{
			description: "CreateUser with no users table",
			action: func(svc *UserService) error {
				_, err := svc.CreateUser("Joe", "Bloggs")
				return err
			},
		},
		{
			description: "GetAllUsers with no users table",
			action: func(svc *UserService) error {
				_, err := svc.GetAllUsers()
				return err
			},
		},
		{
			description: "GetUser with no users table",
			action: func(svc *UserService) error {
				_, err := svc.GetUser(1)
				return err
			},
		},
		{
			description: "UpdateUser with no users table",
			action: func(svc *UserService) error {
				firstName := "Joe"
				lastName := "Bloggs"
				_, err := svc.UpdateUser(1, &firstName, &lastName)
				return err
			},
		},
		{
			description: "DeleteUser with no users table",
			action: func(svc *UserService) error {
				return svc.DeleteUser(1)
			},
		},
	}

	// Test setup / arrangement
	connectionString := "file:user_test_no_tables?mode=memory&cache=shared"
	conn, err := gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("Database failed to connect: " + err.Error())
	}
	db := &database.Database{Conn: conn}
	// Deliberately creating a DB without any tables,
	// i.e. automigrations are not being run:

	// modelsToMigrate := []interface{}{&models.User{}}
	// db.Conn.AutoMigrate(modelsToMigrate...)

	userService := &UserService{DB: db}

	expectedErrorMessage := "no such table: users"
	executeDbTests(userService, t, testCases, &expectedErrorMessage)
}

func executeDbTests(svc *UserService, t *testing.T, testCases []databaseTest, expectedErrorMessage *string) {
	for _, test := range testCases {
		t.Run(fmt.Sprintf("%s - %s", t.Name(), test.description), func(t *testing.T) {
			err := test.action(svc)
			assert.NotNil(t, err, fmt.Sprintf("%s did not return the expected error.", test.description))
			assert.Equal(t, *expectedErrorMessage, err.Error())
		})
	}
}
