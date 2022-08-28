package main

import (
	"log"

	"github.com/conormkelly/fiber-demo/controllers"
	"github.com/conormkelly/fiber-demo/database"
	"github.com/conormkelly/fiber-demo/models"
	"github.com/gofiber/fiber/v2"
)

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
	usersController := &controllers.UsersController{DB: app.DB.Conn}

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
