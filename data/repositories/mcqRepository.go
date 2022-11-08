package repositories

import "gorm.io/gorm"

type McqRepository struct {
	DBConn *gorm.DB
}
