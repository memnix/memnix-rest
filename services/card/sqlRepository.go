package card

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/memnix/memnix-rest/db/sqlc"
	"github.com/memnix/memnix-rest/domain"
)

type SQLRepository struct {
	q *db.Queries // q is the sqlc queries.
}

func NewRepository(dbConn *pgxpool.Pool) IRepository {
	q := db.New(dbConn)
	return &SQLRepository{q: q}
}

func (r SQLRepository) GetByID(_ context.Context, _ uint) (domain.Card, error) {
	return domain.Card{}, nil
}
