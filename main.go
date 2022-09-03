package main

import (
	"log"

	"github.com/conormkelly/fiber-demo/controllers"
	"github.com/conormkelly/fiber-demo/database"
	"github.com/conormkelly/fiber-demo/models"
	"github.com/conormkelly/fiber-demo/services"
	"github.com/gofiber/fiber/v2"
)

type App struct {
	Fiber *fiber.App
	DB    *database.Database
}

func (app *App) Initialize(db *database.Database) {
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
	db := &database.Database{}
	connectionString := "./database/test.db"
	modelsToMigrate := []interface{}{&models.User{}}
	err := db.Connect(&database.Options{SQLitePath: &connectionString, ModelsToMigrate: modelsToMigrate})
	if err != nil {
		log.Fatal("Database failed to connect: " + err.Error())
	}

	app := &App{}
	app.Initialize(db)

	log.Fatal(app.Fiber.Listen(":3000"))
}
