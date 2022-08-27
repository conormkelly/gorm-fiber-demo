package main

import (
	"log"

	"github.com/conormkelly/fiber-demo/database"
	"github.com/conormkelly/fiber-demo/routes"
	"github.com/gofiber/fiber/v2"
)

func initializeRoutes(app *fiber.App) {
	// Users
	app.Post("/api/users", routes.CreateUser)
	app.Get("/api/users", routes.GetAllUsers)
	app.Get("/api/users/:id", routes.GetUserById)
	app.Put("/api/users/:id", routes.UpdateUser)
	app.Delete("/api/users/:id", routes.DeleteUser)
}

func main() {
	db := &database.SQLDatabase{}
	db.Connect()

	app := fiber.New()
	initializeRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
