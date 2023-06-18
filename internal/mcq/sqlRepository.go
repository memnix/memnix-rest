package mcq

import (
	"github.com/memnix/memnix-rest/domain"
	"gorm.io/gorm"
)

type SQLRepository struct {
	DBConn *gorm.DB
}

func (r SQLRepository) GetByID(id uint) (domain.Mcq, error) {
	var mcq domain.Mcq
	err := r.DBConn.First(&mcq, id).Error
	return mcq, err
}

func NewRepository(dbConn *gorm.DB) IRepository {
	return SQLRepository{DBConn: dbConn}
}
