package controllers

import (
	"io"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"github.com/conormkelly/fiber-demo/database"
)

var app App

type App struct {
	Fiber *fiber.App
	DB    *database.Database
}

func (app *App) Initialize(db *database.Database) {
	app.DB = db
	fiberApp := fiber.New()
	app.Fiber = fiberApp

	app.initializeRoutes()
}

func (app *App) initializeRoutes() {
	usersController := &UsersController{DB: app.DB.Conn}

	app.Fiber.Post("/api/users", usersController.CreateUser)
	app.Fiber.Get("/api/users", usersController.GetAllUsers)
	app.Fiber.Get("/api/users/:id", usersController.GetUserById)
	app.Fiber.Put("/api/users/:id", usersController.UpdateUser)
	app.Fiber.Delete("/api/users/:id", usersController.DeleteUser)
}

// Create an in-memory SQLite DB for testing purposes.
func TestMain(m *testing.M) {
	db := database.Database{}
	db.Connect(&database.Options{UseInMemoryDatabase: true, PerformMigration: true})

	app.Initialize(&db)

	code := m.Run()
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	testCases := []testCase{
		{
			description:        "Create valid user",
			method:             "POST",
			route:              "/api/users",
			body:               strings.NewReader(`{ "first_name": "Harry", "last_name": "Potter" }`),
			expectedStatusCode: 200,
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
	description        string    // description of the test case
	method             string    // GET POST etc
	route              string    // route path to test
	body               io.Reader // JSON request body
	expectedStatusCode int       // HTTP status code
	expectedResponse   string    // response body
}

// Helper function
func executeTests(t *testing.T, testCases []testCase) {
	// Iterate through test single test cases
	for _, test := range testCases {
		// Create a new http request with the route from the test case
		req := httptest.NewRequest(test.method, test.route, test.body)
		req.Header.Set("Content-Type", "application/json")

		// Perform the request plain with the app,
		// the second argument is a request latency
		// (set to -1 for no latency)
		resp, _ := app.Fiber.Test(req, -1)

		assert.Equalf(t, test.expectedStatusCode, resp.StatusCode, test.description)

		if test.expectedResponse != "" {
			body, err := ioutil.ReadAll(resp.Body)
			assert.Nil(t, err, "Error parsing response body")
			actualResponse := string(body)
			assert.Equalf(t, test.expectedResponse, actualResponse, test.description)
		}
	}
}
