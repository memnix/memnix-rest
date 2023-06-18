package card

import (
	"context"

	"github.com/memnix/memnix-rest/domain"
	"gorm.io/gorm"
)

type SQLRepository struct {
	DBConn *gorm.DB
}

func NewRepository(dbConn *gorm.DB) IRepository {
	return &SQLRepository{DBConn: dbConn}
}

func (r SQLRepository) GetByID(ctx context.Context, id uint) (domain.Card, error) {
	var card domain.Card
	err := r.DBConn.WithContext(ctx).First(&card, id).Error
	return card, err
}
