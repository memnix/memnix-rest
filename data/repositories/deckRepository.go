package repositories

import "gorm.io/gorm"

type DeckRepository struct {
	DBConn *gorm.DB
}
