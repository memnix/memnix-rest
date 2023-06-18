package card

import (
	"github.com/memnix/memnix-rest/domain"
	"gorm.io/gorm"
)

type SQLRepository struct {
	DBConn *gorm.DB
}

func NewRepository(dbConn *gorm.DB) IRepository {
	return &SQLRepository{DBConn: dbConn}
}

func (r SQLRepository) GetByID(id uint) (domain.Card, error) {
	var card domain.Card
	err := r.DBConn.First(&card, id).Error
	return card, err
}
