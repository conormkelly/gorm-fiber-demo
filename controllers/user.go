package controllers

import (
	"errors"
	"log"

	"github.com/conormkelly/fiber-demo/models"
	"github.com/conormkelly/fiber-demo/services"
	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	Message string `json:"message"`
}

// This is not the user model,
// think of it as a serializer
type User struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func Serialize(userModel models.User) User {
	return User{ID: userModel.ID, FirstName: userModel.FirstName, LastName: userModel.LastName}
}

type UsersController struct {
	Service *services.UserService
}

func ParseBody(ctx *fiber.Ctx, target interface{}) error {
	if err := ctx.BodyParser(target); err != nil {
		return errors.New("invalid JSON request body provided")
	}
	return nil
}

func (c *UsersController) CreateUser(ctx *fiber.Ctx) error {
	user := &models.User{}
	if err := ParseBody(ctx, user); err != nil {
		return ctx.Status(400).JSON(APIResponse{Message: err.Error()})
	}

	user, err := c.Service.CreateUser(user.FirstName, user.LastName)
	if err != nil {
		log.Printf("Error occurred in svc.CreateUser: " + err.Error())
		return err
	}

	serializedUser := Serialize(*user)

	return ctx.Status(200).JSON(serializedUser)
}

func (c *UsersController) GetAllUsers(ctx *fiber.Ctx) error {
	users, err := c.Service.GetAllUsers()
	if err != nil {
		log.Printf("Error occurred in svc.GetAllUsers: " + err.Error())
		return err
	}

	var serializedUsers = make([]User, len(users))
	for i, user := range users {
		serializedUsers[i] = Serialize(user)
	}

	return ctx.Status(200).JSON(serializedUsers)
}

func (c *UsersController) GetUserById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(400).JSON(APIResponse{Message: "User ID must be an integer"})
	}

	user, err := c.Service.GetUser(id)
	if err != nil {
		return err
	}

	serializedUser := Serialize(*user)
	return ctx.Status(200).JSON(serializedUser)
}

func (c *UsersController) UpdateUser(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(400).JSON(APIResponse{Message: "User ID must be an integer"})
	}

	var updatedUser models.User
	if err := ParseBody(ctx, &updatedUser); err != nil {
		return ctx.Status(400).JSON(APIResponse{Message: err.Error()})
	}

	user, err := c.Service.UpdateUser(id, &updatedUser.FirstName, &updatedUser.LastName)
	if err != nil {
		return err
	}

	serializedUser := Serialize(*user)
	return ctx.Status(200).JSON(serializedUser)
}

func (c *UsersController) DeleteUser(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(400).JSON(APIResponse{Message: "User ID must be an integer"})
	}

	err = c.Service.DeleteUser(id)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(APIResponse{Message: "Successfully deleted user"})
}
