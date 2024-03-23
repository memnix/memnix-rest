package deck

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

// Create creates a new deck.
func (r *SQLRepository) Create(_ context.Context, _ *domain.Deck) error {
	return nil
}

// Update updates the deck with the given id.
func (r *SQLRepository) Update(_ context.Context, _ *domain.Deck) error {
	return nil
}

// Delete deletes the deck with the given id.
func (r *SQLRepository) Delete(_ context.Context, _ uint) error {
	return nil
}

// CreateFromUser creates a new deck from the given user.
func (r *SQLRepository) CreateFromUser(_ context.Context, _ domain.User, _ *domain.Deck) error {
	return nil
}

// GetByID returns the deck with the given id.
func (r *SQLRepository) GetByID(_ context.Context, _ uint) (domain.Deck, error) {
	return domain.Deck{}, nil
}

// GetByUser returns the decks of the given user.
func (r *SQLRepository) GetByUser(_ context.Context, _ domain.User) ([]domain.Deck, error) {
	return []domain.Deck{}, nil
}

// GetByLearner returns the decks of the given learner.
func (r *SQLRepository) GetByLearner(_ context.Context, _ domain.User) ([]domain.Deck, error) {
	return []domain.Deck{}, nil
}

// GetPublic returns the public decks.
func (r *SQLRepository) GetPublic(_ context.Context) ([]domain.Deck, error) {
	return []domain.Deck{}, nil
}

// NewRepository returns a new repository.
func NewRepository(dbConn *pgxpool.Pool) IRepository {
	q := db.New(dbConn)
	return &SQLRepository{q: q}
}
