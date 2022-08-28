package controllers

import (
	"errors"

	"github.com/conormkelly/fiber-demo/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type APIResponse struct {
	Message string `json:"message"`
}

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

func (controller *UsersController) CreateUser(c *fiber.Ctx) error {
	user := &models.User{}

	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(APIResponse{Message: err.Error()})
	}

	controller.DB.Create(&user)
	serializedUser := Serialize(*user)

	return c.Status(200).JSON(serializedUser)
}

func (controller *UsersController) GetAllUsers(c *fiber.Ctx) error {
	users := []models.User{}
	controller.DB.Find(&users)

	var serializedUsers = make([]User, len(users))
	for i, user := range users {
		serializedUsers[i] = Serialize(user)
	}

	return c.Status(200).JSON(serializedUsers)
}

// Helper function, TOOD: refactor
func (controller *UsersController) findUser(id int, user *models.User) error {
	controller.DB.Find(&user, "id = ?", id)
	if user.ID == 0 {
		return errors.New("User does not exist")
	}
	return nil
}

func (controller *UsersController) GetUserById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(APIResponse{Message: "User ID must be an integer"})
	}

	var user models.User
	if err := controller.findUser(id, &user); err != nil {
		return c.Status(400).JSON(APIResponse{Message: err.Error()})
	}

	serializedUser := Serialize(user)
	return c.Status(200).JSON(serializedUser)
}

func (controller *UsersController) UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(APIResponse{Message: "User ID must be an integer"})
	}

	var user models.User
	if err := controller.findUser(id, &user); err != nil {
		return c.Status(400).JSON(APIResponse{Message: err.Error()})
	}

	var updatedUser models.User
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(400).JSON(APIResponse{Message: err.Error()})
	}

	if updatedUser.FirstName != "" {
		user.FirstName = updatedUser.FirstName
	}
	if updatedUser.LastName != "" {
		user.LastName = updatedUser.LastName
	}

	controller.DB.Save(user)

	serializedUser := Serialize(user)
	return c.Status(200).JSON(serializedUser)
}

func (controller *UsersController) DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(APIResponse{Message: "User ID must be an integer"})
	}

	var user models.User
	if err := controller.findUser(id, &user); err != nil {
		return c.Status(400).JSON(APIResponse{Message: err.Error()})
	}

	if err := controller.DB.Delete(user).Error; err != nil {
		return c.Status(500).JSON(APIResponse{Message: err.Error()})
	}

	return c.Status(200).JSON(APIResponse{Message: "Successfully deleted user"})
}
