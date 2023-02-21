package user

import (
	"github.com/memnix/memnix-rest/domain"
	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
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
	log.Info().Msgf("user: %v", user)
	return user.Username
}
