package main

import (
	"errors"
	"log"
	"os"

	"github.com/conormkelly/fiber-demo/controllers"
	"github.com/conormkelly/fiber-demo/database"
	"github.com/conormkelly/fiber-demo/models"
	"github.com/conormkelly/fiber-demo/services"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-multierror"
)

type Options struct {
	DatabaseType     database.DatabaseType // SQLite, MySQL etc. Only SQLite is currently supported.
	ConnectionString *string               // The path to the DB e.g. "file::memory:?cache=shared"
	ModelsToMigrate  []interface{}         // (optional) The models that should be migrated
	Port             *string               // (optional) If not supplied, the app wont actually listen
}

type App struct {
	Options *Options
	Fiber   *fiber.App
	DB      *database.Database
}

// Creates DB connection based on supplied config,
// sets up routes and listens if the port was supplied.
func (app *App) Initialize(configError error) error {
	if configError != nil {
		return configError
	}

	db := &database.Database{}

	dbOptions := &database.Options{
		DatabaseType:     app.Options.DatabaseType,
		ConnectionString: app.Options.ConnectionString,
		ModelsToMigrate:  app.Options.ModelsToMigrate,
	}

	err := db.Connect(dbOptions)
	if err != nil {
		return err
	}
	app.DB = db

	// Create a new fiber instance with custom config
	fiberApp := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's an fiber.*Error
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			} else {
				// Log, but return a generic error to client to avoid leaking error details
				log.Printf("An application error occured: " + err.Error())
				err = fiber.NewError(500, "sorry, something went wrong")
			}

			// Send custom error
			return ctx.Status(code).JSON(controllers.APIResponse{Message: err.Error()})
		},
	})

	app.Fiber = fiberApp

	app.initializeRoutes()

	if app.Options.Port != nil {
		log.Fatal(app.Fiber.Listen(*app.Options.Port))
	}

	return nil
}

func (app *App) initializeRoutes() {
	usersController := &controllers.UsersController{Service: &services.UserService{DB: app.DB}}

	app.Fiber.Post("/api/users", usersController.CreateUser)
	app.Fiber.Get("/api/users", usersController.GetAllUsers)
	app.Fiber.Get("/api/users/:id", usersController.GetUserById)
	app.Fiber.Put("/api/users/:id", usersController.UpdateUser)
	app.Fiber.Delete("/api/users/:id", usersController.DeleteUser)
}

// Reads and validate environment variable config
func ReadEnvironmentConfig() (*Options, error) {
	options := Options{}
	var result error

	// Determine DB type
	dbTypeString := os.Getenv("APP_DB_TYPE")
	dbType := database.Undefined

	switch {
	case dbTypeString == "":
		err := errors.New(`APP_DB_TYPE should be one of: ["MYSQL", "SQLITE"]`)
		result = multierror.Append(result, err)
	case dbTypeString == "MYSQL":
		dbType = database.MySQL
	case dbTypeString == "SQLITE":
		dbType = database.SQLite
	}
	options.DatabaseType = dbType

	// Get connection string
	dbConnectionString := os.Getenv("APP_DB_CONN_STRING")
	if dbConnectionString == "" {
		err := errors.New("APP_DB_CONN_STRING is required")
		result = multierror.Append(result, err)
		options.ConnectionString = nil
	} else {
		options.ConnectionString = &dbConnectionString
	}

	// Get port
	port := os.Getenv("APP_PORT")
	if port == "" {
		err := errors.New("APP_PORT is required")
		result = multierror.Append(result, err)
		options.Port = nil
	} else {
		options.Port = &port
	}

	// Check if we should run migrations
	shouldRunMigrations := os.Getenv("APP_RUN_AUTO_MIGRATE")
	if shouldRunMigrations == "true" {
		options.ModelsToMigrate = []interface{}{&models.User{}}
	}

	return &options, result
}

func main() {
	options, configError := ReadEnvironmentConfig()
	app := &App{Options: options}
	err := app.Initialize(configError)
	if err != nil {
		log.Fatal("Failed to start app: " + err.Error())
	}
}
