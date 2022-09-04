package main

import (
	"log"

	"github.com/conormkelly/fiber-demo/controllers"
	"github.com/conormkelly/fiber-demo/database"
	"github.com/conormkelly/fiber-demo/models"
	"github.com/conormkelly/fiber-demo/services"
	"github.com/gofiber/fiber/v2"
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
func (app *App) Initialize() error {
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

func main() {
	// TODO: get actual config from env / Viper
	dbType := database.SQLite
	connectionString := "./database/test.db"
	port := ":3000"

	options := Options{
		DatabaseType:     dbType,
		ConnectionString: &connectionString,
		ModelsToMigrate:  []interface{}{&models.User{}},
		Port:             &port,
	}

	app := &App{Options: &options}

	err := app.Initialize()
	if err != nil {
		log.Fatal("Failed to init app: " + err.Error())
	}
}
