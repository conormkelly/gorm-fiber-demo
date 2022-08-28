package services

import (
	"github.com/conormkelly/fiber-demo/database"
	"github.com/conormkelly/fiber-demo/models"
	"github.com/gofiber/fiber/v2"
)

type UserService struct {
	DB *database.Database
}

func (svc *UserService) CreateUser(firstName, lastName string) (*models.User, error) {
	user := &models.User{FirstName: firstName, LastName: lastName}

	err := svc.DB.Conn.Create(&user).Error
	return user, err
}

func (svc *UserService) GetAllUsers() ([]models.User, error) {
	users := []models.User{}
	err := svc.DB.Conn.Find(&users).Error

	return users, err
}

func (svc *UserService) GetUser(id int) (*models.User, error) {
	var user models.User
	err := svc.DB.Conn.Find(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	} else if user.ID == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "user does not exist")
	}
	return &user, nil
}

func (svc *UserService) UpdateUser(id int, firstName, lastName *string) (*models.User, error) {
	user, err := svc.GetUser(id)
	if err != nil {
		return nil, err
	}
	if firstName != nil && *firstName != "" {
		user.FirstName = *firstName
	}
	if lastName != nil && *lastName != "" {
		user.LastName = *lastName
	}

	err = svc.DB.Conn.Save(user).Error

	return user, err
}

func (svc *UserService) DeleteUser(id int) error {
	user, err := svc.GetUser(id)
	if err != nil {
		return err
	}

	return svc.DB.Conn.Delete(user).Error
}
