package mcq

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/memnix/memnix-rest/db/sqlc"
	"github.com/memnix/memnix-rest/domain"
)

// SQLRepository is the repository for the deck.
type SQLRepository struct {
	q *db.Queries // q is the sqlc queries.
}

func (r SQLRepository) GetByID(_ context.Context, _ uint) (domain.Mcq, error) {
	return domain.Mcq{}, nil
}

// NewRepository returns a new repository.
func NewRepository(dbConn *pgxpool.Pool) IRepository {
	q := db.New(dbConn)
	return &SQLRepository{q: q}
}
