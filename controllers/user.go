package controllers

import (
	"errors"

	"github.com/conormkelly/fiber-demo/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func Serialize(userModel models.User) User {
	return User{ID: userModel.ID, FirstName: userModel.FirstName, LastName: userModel.LastName}
}

type UsersController struct {
	DB *gorm.DB
}

func (handler *UsersController) CreateUser(c *fiber.Ctx) error {
	user := &models.User{}

	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	handler.DB.Create(&user)
	serializedUser := Serialize(*user)

	return c.Status(200).JSON(serializedUser)
}

func (handler *UsersController) GetAllUsers(c *fiber.Ctx) error {
	users := []models.User{}
	handler.DB.Find(&users)

	var serializedUsers = make([]User, len(users))
	for i, user := range users {
		serializedUsers[i] = Serialize(user)
	}

	return c.Status(200).JSON(serializedUsers)
}

// Helper function, TOOD: refactor
func (handler *UsersController) findUser(id int, user *models.User) error {
	handler.DB.Find(&user, "id = ?", id)
	if user.ID == 0 {
		return errors.New("User does not exist")
	}
	return nil
}

func (handler *UsersController) GetUserById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("User ID must be an integer")
	}

	var user models.User
	if err := handler.findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	serializedUser := Serialize(user)
	return c.Status(200).JSON(serializedUser)
}

func (handler *UsersController) UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("User ID must be an integer")
	}

	var user models.User
	if err := handler.findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var updatedUser models.User
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if updatedUser.FirstName != "" {
		user.FirstName = updatedUser.FirstName
	}
	if updatedUser.LastName != "" {
		user.LastName = updatedUser.LastName
	}

	handler.DB.Save(user)

	serializedUser := Serialize(user)
	return c.Status(200).JSON(serializedUser)
}

func (handler *UsersController) DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("User ID must be an integer")
	}

	var user models.User
	if err := handler.findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := handler.DB.Delete(user).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON("Successfully deleted user")
}
