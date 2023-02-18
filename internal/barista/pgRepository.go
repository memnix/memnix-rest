package barista

import "gorm.io/gorm"

type PgRepository struct {
	DBConn *gorm.DB
}

func NewPgRepository(dbConn *gorm.DB) IPgRepository {
	return &PgRepository{DBConn: dbConn}
}

// GetName returns the name of the repository.
func (r *PgRepository) GetName() string {
	return "Sexy barista"
}
