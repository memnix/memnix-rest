package user

import (
	"github.com/memnix/memnix-rest/domain"
	"gorm.io/gorm"
)

type SqlRepository struct {
	DBConn *gorm.DB
}

func NewRepository(dbConn *gorm.DB) IRepository {
	return &SqlRepository{DBConn: dbConn}
}

// GetName returns the name of the user.
func (r *SqlRepository) GetName(id uint) string {
	var user domain.User
	r.DBConn.First(&user, id)
	return user.Username
}

// GetByID returns the user with the given id.
func (r *SqlRepository) GetByID(id uint) (domain.User, error) {
	var user domain.User
	err := r.DBConn.First(&user, id).Error
	return user, err
}

// GetByEmail returns the user with the given email.
func (r *SqlRepository) GetByEmail(email string) (domain.User, error) {
	var user domain.User
	err := r.DBConn.Where("email = ?", email).First(&user).Error
	return user, err
}

// Create creates a new user.
func (r *SqlRepository) Create(user *domain.User) error {
	return r.DBConn.Create(&user).Error
}

func (r *SqlRepository) Update(user *domain.User) error {
	return r.DBConn.Save(&user).Error
}

func (r *SqlRepository) Delete(id uint) error {
	return r.DBConn.Delete(&domain.User{}, id).Error
}

func (r *SqlRepository) GetByOauthID(id string) (domain.User, error) {
	var user domain.User
	err := r.DBConn.Where("oauth_id = ?", id).First(&user).Error
	return user, err
}
