package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/conormkelly/fiber-demo/database"
	"github.com/conormkelly/fiber-demo/models"
)

var app App

// Create an in-memory SQLite DB for testing purposes.
func TestMain(m *testing.M) {
	db := database.Database{}
	modelsToMigrate := []interface{}{&models.User{}}
	err := db.Connect(&database.Options{UseInMemoryDatabase: true, InMemoryDatabaseName: "main_app", ModelsToMigrate: modelsToMigrate})
	if err != nil {
		log.Fatal("Database failed to connect: " + err.Error())
	}

	app.Initialize(&db)

	code := m.Run()
	os.Exit(code)
}

func Test404Handler(t *testing.T) {
	test := testCase{
		description:        "Test non-existent route",
		method:             "GET",
		route:              "/api/non-existent-route",
		expectedStatusCode: 404,
		expectedResponse:   `{"message":"Cannot GET /api/non-existent-route"}`,
	}
	executeTest(t, &app, test)
}

func TestCreateUser(t *testing.T) {
	testCases := []testCase{
		{
			description:        "Create valid user",
			method:             "POST",
			route:              "/api/users",
			body:               strings.NewReader(`{ "first_name": "John", "last_name": "Doe" }`),
			expectedStatusCode: 200,
		},
		{
			description:        "Malformed JSON",
			method:             "POST",
			route:              "/api/users",
			body:               strings.NewReader(`{ "first_name": NO CLOSING BRACKET`),
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"invalid JSON request body provided"}`,
		},
		// {
		// 	description:  "Create partial user",
		// 	method:       "POST",
		// 	route:        "/api/users",
		// 	body:         strings.NewReader(`{ "first_name": "John" }`),
		// 	expectedCode: 400,
		// },
		{
			description:        "Send invalid JSON",
			method:             "POST",
			route:              "/api/users",
			body:               strings.NewReader(`{ INVALID JSON `),
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"invalid JSON request body provided"}`,
		},
		{
			description:        "Send no JSON",
			method:             "POST",
			route:              "/api/users",
			body:               nil,
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"invalid JSON request body provided"}`,
		},
	}

	executeTests(t, &app, testCases)
}
func TestGetAllUsers(t *testing.T) {
	testCases := []testCase{
		{
			description:        "Get all users when table not empty",
			method:             "GET",
			route:              "/api/users",
			expectedStatusCode: 200,
			expectedResponse:   `[{"id":1,"first_name":"John","last_name":"Doe"}]`,
			setup: func() {
				clearTable(&app)
				addUser(&app)
			},
		},
		{
			description:        "Get all users when table is empty",
			method:             "GET",
			route:              "/api/users",
			expectedStatusCode: 200,
			expectedResponse:   `[]`,
			setup: func() {
				clearTable(&app)
			},
		},
	}

	executeTests(t, &app, testCases)
}

func TestGetUserById(t *testing.T) {
	testCases := []testCase{
		{
			description:        "Get user by ID",
			method:             "GET",
			route:              "/api/users/1",
			expectedStatusCode: 200,
			expectedResponse:   `{"id":1,"first_name":"John","last_name":"Doe"}`,
			setup: func() {
				clearTable(&app)
				addUser(&app)
			},
		},
		{
			description:        "Get user by non-integer ID",
			method:             "GET",
			route:              "/api/users/one",
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"User ID must be an integer"}`,
		},
		{
			description:        "Get non-existent user",
			method:             "GET",
			route:              "/api/users/1",
			expectedStatusCode: 404,
			expectedResponse:   `{"message":"user does not exist"}`,
			setup: func() {
				clearTable(&app)
			},
		},
	}

	executeTests(t, &app, testCases)
}

func TestUpdateUser(t *testing.T) {
	testCases := []testCase{
		{
			description:        "Update existing user",
			method:             "PUT",
			route:              "/api/users/1",
			body:               strings.NewReader(`{"first_name":"James","last_name":"Doe"}`),
			expectedStatusCode: 200,
			expectedResponse:   `{"id":1,"first_name":"James","last_name":"Doe"}`,
			setup: func() {
				clearTable(&app)
				addUser(&app)
			},
		},
		{
			description:        "Update user with non-int id",
			method:             "PUT",
			route:              "/api/users/two",
			body:               strings.NewReader(`{"first_name":"James","last_name":"Doe"}`),
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"User ID must be an integer"}`,
		},
		{
			description:        "Malformed JSON",
			method:             "PUT",
			route:              "/api/users/1",
			body:               strings.NewReader(`{ "first_name": NO CLOSING BRACKET`),
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"invalid JSON request body provided"}`,
		},
		{
			description:        "Update non-existent user",
			method:             "PUT",
			route:              "/api/users/0",
			body:               strings.NewReader(`{"first_name":"James","last_name":"Doe"}`),
			expectedStatusCode: 404,
			expectedResponse:   `{"message":"user does not exist"}`,
		},
	}

	executeTests(t, &app, testCases)
}

func TestDeleteUser(t *testing.T) {
	testCases := []testCase{
		{
			description:        "Delete existing user",
			method:             "DELETE",
			route:              "/api/users/1",
			expectedStatusCode: 200,
			expectedResponse:   `{"message":"Successfully deleted user"}`,
			setup: func() {
				clearTable(&app)
				addUser(&app)
			},
		},
		{
			description:        "Delete non-existent user",
			method:             "DELETE",
			route:              "/api/users/1",
			expectedStatusCode: 404,
			expectedResponse:   `{"message":"user does not exist"}`,
			setup: func() {
				clearTable(&app)
			},
		},
		{
			description:        "Delete user with non-int id",
			method:             "DELETE",
			route:              "/api/users/three",
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"User ID must be an integer"}`,
		},
	}

	executeTests(t, &app, testCases)
}

// Check that API returns correctly sanitized error messages when DB is not in good state
func TestDBErrors(t *testing.T) {
	// Test setup / arrangement

	// An app with no tables migrated
	var brokenApp App
	db := database.Database{}
	err := db.Connect(&database.Options{UseInMemoryDatabase: true, InMemoryDatabaseName: "broken_app"})
	if err != nil {
		log.Fatal("Database failed to connect: " + err.Error())
	}

	brokenApp.Initialize(&db)

	testCases := []testCase{
		{
			description:        "Create user with no users table",
			method:             "POST",
			route:              "/api/users",
			body:               strings.NewReader(`{ "first_name": "John", "last_name": "Doe" }`),
			expectedStatusCode: 500,
			expectedResponse:   `{"message":"sorry, something went wrong"}`,
		},
		{
			description:        "Get all users with no users table",
			method:             "GET",
			route:              "/api/users",
			expectedStatusCode: 500,
			expectedResponse:   `{"message":"sorry, something went wrong"}`,
		},
		{
			description:        "Get user by ID with no users table",
			method:             "GET",
			route:              "/api/users/1",
			expectedStatusCode: 500,
			expectedResponse:   `{"message":"sorry, something went wrong"}`,
		},
		{
			description:        "Update user with no users table",
			method:             "PUT",
			route:              "/api/users/1",
			body:               strings.NewReader(`{ "last_name": "Bond", "first_name": "James" }`),
			expectedStatusCode: 500,
			expectedResponse:   `{"message":"sorry, something went wrong"}`,
		},
		{
			description:        "Delete user with no users table",
			method:             "DELETE",
			route:              "/api/users/1",
			expectedStatusCode: 500,
			expectedResponse:   `{"message":"sorry, something went wrong"}`,
		},
	}

	executeTests(t, &brokenApp, testCases)
}

// Helper methods

type testCase struct {
	description        string    // Description of the test case
	method             string    // GET, POST etc
	route              string    // Endpoint to test
	body               io.Reader // JSON request body
	expectedStatusCode int       // HTTP status code
	expectedResponse   string    // Response body
	setup              func()    // Hook to run a function before the test executes
}

func executeTest(t *testing.T, application *App, test testCase) {
	t.Run(fmt.Sprintf("%s - %s", t.Name(), test.description), func(t *testing.T) {
		// Hook to run a function before the test executes
		if test.setup != nil {
			test.setup()
		}
		// Create a new HTTP request with the route from the test case
		req := httptest.NewRequest(test.method, test.route, test.body)
		req.Header.Set("Content-Type", "application/json")

		// Perform the request against the Fiber app,
		// with a timeout of 500ms
		resp, err := application.Fiber.Test(req, 500)
		assert.Nil(t, err, "Fiber.Test returned an error")

		assert.Equalf(t, test.expectedStatusCode, resp.StatusCode, test.description)

		if test.expectedResponse != "" {
			body, err := ioutil.ReadAll(resp.Body)
			assert.Nil(t, err, "Error parsing response body")
			actualResponse := string(body)
			assert.Equalf(t, test.expectedResponse, actualResponse, test.description)
		}
	})
}

func executeTests(t *testing.T, application *App, testCases []testCase) {
	// Iterate through test cases
	for _, test := range testCases {
		executeTest(t, application, test)
	}
}

func clearTable(application *App) {
	application.DB.Conn.Where("id > ?", 0).Delete(&models.User{})
}

func addUser(application *App) {
	user := &models.User{FirstName: "John", LastName: "Doe"}
	application.DB.Conn.Create(user)
}
