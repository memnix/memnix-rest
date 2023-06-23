package mcq

import (
	"context"

	"github.com/memnix/memnix-rest/domain"
	"gorm.io/gorm"
)

type SQLRepository struct {
	DBConn *gorm.DB
}

func (r SQLRepository) GetByID(ctx context.Context, id uint) (domain.Mcq, error) {
	var mcq domain.Mcq
	err := r.DBConn.WithContext(ctx).First(&mcq, id).Error
	return mcq, err
}

func NewRepository(dbConn *gorm.DB) IRepository {
	return SQLRepository{DBConn: dbConn}
}
