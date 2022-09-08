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
	ConnectionString  *string // The DSN for the DB e.g. "user:password@tcp(localhost:3306)/go_app"
	ShouldAutoMigrate bool
	Port              *string
}

type App struct {
	Options *Options
	Fiber   *fiber.App
	DB      *database.Database
}

// Parse environment variable config into Options
func GetAppOptions() (*Options, error) {
	options := Options{}
	var configErrors error

	dbConnectionString := os.Getenv("APP_DB_CONN_STRING")
	if dbConnectionString == "" {
		err := errors.New("APP_DB_CONN_STRING is required")
		configErrors = multierror.Append(configErrors, err)
		options.ConnectionString = nil
	} else {
		options.ConnectionString = &dbConnectionString
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		err := errors.New("APP_PORT is required")
		configErrors = multierror.Append(configErrors, err)
		options.Port = nil
	} else {
		options.Port = &port
	}

	options.ShouldAutoMigrate = os.Getenv("APP_RUN_AUTO_MIGRATE") == "true"

	return &options, configErrors
}

// Creates DB connection based on supplied config
func (app *App) ConnectDB() error {
	var modelsToMigrate = []interface{}{}
	if app.Options.ShouldAutoMigrate {
		modelsToMigrate = []interface{}{&models.User{}}
	}

	dbOptions := &database.Options{
		ConnectionString: app.Options.ConnectionString,
		ModelsToMigrate:  modelsToMigrate,
	}

	conn, err := database.GetConnection(dbOptions)
	if err != nil {
		return err
	}
	app.DB = &database.Database{Conn: conn}
	return nil
}

func (app *App) ConfigureFiber() {
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
}

func (app *App) InitializeRoutes() {
	usersController := &controllers.UsersController{Service: &services.UserService{DB: app.DB}}

	app.Fiber.Post("/api/users", usersController.CreateUser)
	app.Fiber.Get("/api/users", usersController.GetAllUsers)
	app.Fiber.Get("/api/users/:id", usersController.GetUserById)
	app.Fiber.Put("/api/users/:id", usersController.UpdateUser)
	app.Fiber.Delete("/api/users/:id", usersController.DeleteUser)

	if app.Options.Port != nil {
		log.Fatal(app.Fiber.Listen(*app.Options.Port))
	}
}

func Start() error {
	// Read and validate config
	options, configError := GetAppOptions()
	if configError != nil {
		return errors.New("invalid config - " + configError.Error())
	}

	// Create app and connect to DB
	app := &App{Options: options}
	dbError := app.ConnectDB()
	if dbError != nil {
		return errors.New("DB connection error - " + dbError.Error())
	}

	app.ConfigureFiber()
	app.InitializeRoutes()

	return nil
}

func main() {
	err := Start()
	if err != nil {
		log.Fatal("Startup failure: " + err.Error())
	}
}
