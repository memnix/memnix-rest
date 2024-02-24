package user

import (
	"context"

	"github.com/memnix/memnix-rest/domain"
	"gorm.io/gorm"
)

// SQLRepository is the repository for the user.
type SQLRepository struct {
	DBConn *gorm.DB // DBConn is the database connection.
}

// NewRepository returns a new repository.
func NewRepository(dbConn *gorm.DB) IRepository {
	return &SQLRepository{DBConn: dbConn}
}

// GetName returns the name of the user.
func (r *SQLRepository) GetName(ctx context.Context, id uint) string {
	var user domain.User
	r.DBConn.WithContext(ctx).First(&user, id)
	return user.Username
}

// GetByID returns the user with the given id.
func (r *SQLRepository) GetByID(ctx context.Context, id uint) (domain.User, error) {
	var user domain.User
	err := r.DBConn.WithContext(ctx).First(&user, id).Error
	return user, err
}

// GetByEmail returns the user with the given email.
func (r *SQLRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := r.DBConn.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

// Create creates a new user.
func (r *SQLRepository) Create(ctx context.Context, user *domain.User) error {
	return r.DBConn.WithContext(ctx).Create(&user).Error
}

// Update updates the user with the given id.
func (r *SQLRepository) Update(ctx context.Context, user *domain.User) error {
	return r.DBConn.WithContext(ctx).Save(&user).Error
}

// Delete deletes the user with the given id.
func (r *SQLRepository) Delete(ctx context.Context, id uint) error {
	return r.DBConn.WithContext(ctx).Delete(&domain.User{}, id).Error
}

// GetByOauthID returns the user with the given oauth id.
func (r *SQLRepository) GetByOauthID(ctx context.Context, id string) (domain.User, error) {
	var user domain.User
	err := r.DBConn.WithContext(ctx).Where("oauth_id = ?", id).First(&user).Error
	return user, err
}
