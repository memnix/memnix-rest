package repositories

import "gorm.io/gorm"

type AuthRepository struct {
	DBConn *gorm.DB
}
