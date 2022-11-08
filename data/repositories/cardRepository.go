package repositories

import "gorm.io/gorm"

type CardRepository struct {
	DBConn *gorm.DB
}
