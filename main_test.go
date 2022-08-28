package main

import (
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
	err := db.Connect(&database.Options{UseInMemoryDatabase: true, ModelsToMigrate: modelsToMigrate})
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
		body:               nil,
		expectedStatusCode: 404,
		expectedResponse:   `{"message":"Cannot GET /api/non-existent-route"}`,
	}
	executeTest(t, test)
}

func TestCreateUser(t *testing.T) {
	testCases := []testCase{
		{
			description:        "Create valid user",
			method:             "POST",
			route:              "/api/users",
			body:               strings.NewReader(`{ "first_name": "Harry", "last_name": "Potter" }`),
			expectedStatusCode: 200,
			// setup: func() {
			// 	fmt.Println("Example")
			// },
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
		// 	body:         strings.NewReader(`{ "first_name": "Harry" }`),
		// 	expectedCode: 400,
		// },
		{
			description:        "Send invalid JSON",
			method:             "POST",
			route:              "/api/users",
			body:               strings.NewReader(`{ INVALID JSON `),
			expectedStatusCode: 400,
			// TODO: improve response message here
		},
		{
			description:        "Send no JSON",
			method:             "POST",
			route:              "/api/users",
			body:               nil,
			expectedStatusCode: 400,
		},
	}

	executeTests(t, testCases)
}
func TestGetAllUsers(t *testing.T) {
	testCases := []testCase{
		{
			description:        "Get all users",
			method:             "GET",
			route:              "/api/users",
			body:               nil,
			expectedStatusCode: 200,
			expectedResponse:   `[{"id":1,"first_name":"Harry","last_name":"Potter"}]`,
		},
		// TODO: add ability to clear table then re-test
	}

	executeTests(t, testCases)
}

func TestGetUserById(t *testing.T) {
	testCases := []testCase{
		{
			description:        "Get user by ID",
			method:             "GET",
			route:              "/api/users/1",
			body:               nil,
			expectedStatusCode: 200,
			expectedResponse:   `{"id":1,"first_name":"Harry","last_name":"Potter"}`,
		},
		{
			description:        "Get user by non-integer ID",
			method:             "GET",
			route:              "/api/users/one",
			body:               nil,
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"User ID must be an integer"}`,
		},
		{
			description:        "Get user by zero ID",
			method:             "GET",
			route:              "/api/users/0",
			body:               nil,
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"User does not exist"}`,
		},
	}

	executeTests(t, testCases)
}

func TestUpdateUser(t *testing.T) {
	testCases := []testCase{
		{
			description:        "Update existing user",
			method:             "PUT",
			route:              "/api/users/1",
			body:               strings.NewReader(`{"first_name":"Larry","last_name":"Rotter"}`),
			expectedStatusCode: 200,
			expectedResponse:   `{"id":1,"first_name":"Larry","last_name":"Rotter"}`,
		},
		{
			description:        "Malformed JSON",
			method:             "PUT",
			route:              "/api/users/1",
			body:               strings.NewReader(`{ "first_name": NO CLOSING BRACKET`),
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"invalid JSON request body provided"}`,
		},
		// TODO: update non-existent user etc
	}

	executeTests(t, testCases)
}

func TestDeleteUser(t *testing.T) {
	testCases := []testCase{
		{
			description:        "Delete existing user",
			method:             "DELETE",
			route:              "/api/users/1",
			body:               strings.NewReader(`{"first_name":"Larry","last_name":"Rotter"}`),
			expectedStatusCode: 200,
			expectedResponse:   `{"message":"Successfully deleted user"}`,
		},
		{
			description:        "Delete non-existent user",
			method:             "DELETE",
			route:              "/api/users/1",
			body:               strings.NewReader(`{"first_name":"Larry","last_name":"Rotter"}`),
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"User does not exist"}`,
		},
		// TODO: more test cases, improve status code handling
		// also, make tests less temporally dependent
	}

	executeTests(t, testCases)
}

type testCase struct {
	description        string    // Description of the test case
	method             string    // GET, POST etc
	route              string    // Endpoint to test
	body               io.Reader // JSON request body
	expectedStatusCode int       // HTTP status code
	expectedResponse   string    // Response body
	setup              func()    // Hook to run a function before the test executes
}

// Helper functions
func executeTest(t *testing.T, test testCase) {
	// Hook to run a function before the test executes
	if test.setup != nil {
		test.setup()
	}
	// Create a new HTTP request with the route from the test case
	req := httptest.NewRequest(test.method, test.route, test.body)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request against the Fiber app,
	// with a timeout of 500ms
	resp, _ := app.Fiber.Test(req, 500)

	assert.Equalf(t, test.expectedStatusCode, resp.StatusCode, test.description)

	if test.expectedResponse != "" {
		body, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err, "Error parsing response body")
		actualResponse := string(body)
		assert.Equalf(t, test.expectedResponse, actualResponse, test.description)
	}
}

func executeTests(t *testing.T, testCases []testCase) {
	// Iterate through test cases
	for _, test := range testCases {
		executeTest(t, test)
	}
}
