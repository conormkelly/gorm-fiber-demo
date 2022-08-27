package routes

import (
	"errors"

	"github.com/conormkelly/fiber-demo/database"
	"github.com/conormkelly/fiber-demo/models"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func Serialize(userModel models.User) User {
	return User{ID: userModel.ID, FirstName: userModel.FirstName, LastName: userModel.LastName}
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.Database.DB.Create(&user)
	serializedUser := Serialize(user)

	return c.Status(200).JSON(serializedUser)
}

func GetAllUsers(c *fiber.Ctx) error {
	users := []models.User{}
	database.Database.DB.Find(&users)

	var serializedUsers = make([]User, len(users))
	for i, user := range users {
		serializedUsers[i] = Serialize(user)
	}

	return c.Status(200).JSON(serializedUsers)
}

func findUser(id int, user *models.User) error {
	database.Database.DB.Find(&user, "id = ?", id)
	if user.ID == 0 {
		return errors.New("User does not exist")
	}
	return nil
}

func GetUserById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("User ID must be an integer")
	}

	var user models.User
	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	serializedUser := Serialize(user)
	return c.Status(200).JSON(serializedUser)
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("User ID must be an integer")
	}

	var user models.User
	if err := findUser(id, &user); err != nil {
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

	database.Database.DB.Save(user)

	serializedUser := Serialize(user)
	return c.Status(200).JSON(serializedUser)
}

func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("User ID must be an integer")
	}

	var user models.User
	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := database.Database.DB.Delete(user).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON("Successfully deleted user")
}
