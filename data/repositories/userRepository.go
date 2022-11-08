package repositories

import (
	"errors"
	"github.com/memnix/memnixrest/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DBConn *gorm.DB
}

// GetByID returns a user by id
func (u *UserRepository) GetByID(id uint) (models.User, error) {
	var user models.User
	u.DBConn.First(&user, id)
	if user.ID == 0 {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

// GetAll returns all users
func (u *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	u.DBConn.Find(&users)
	return users, nil
}

func (u *UserRepository) Update(user *models.User) error {
	u.DBConn.Save(user)
	return nil
}
