package deck

import (
	"context"

	"github.com/memnix/memnix-rest/domain"
	"gorm.io/gorm"
)

// SQLRepository is the repository for the deck.
type SQLRepository struct {
	DBConn *gorm.DB
}

// Create creates a new deck.
func (r *SQLRepository) Create(ctx context.Context, deck *domain.Deck) error {
	return r.DBConn.WithContext(ctx).Create(&deck).Error
}

// Update updates the deck with the given id.
func (r *SQLRepository) Update(ctx context.Context, deck *domain.Deck) error {
	return r.DBConn.WithContext(ctx).Save(&deck).Error
}

// Delete deletes the deck with the given id.
func (r *SQLRepository) Delete(ctx context.Context, id uint) error {
	return r.DBConn.WithContext(ctx).Delete(&domain.Deck{}, id).Error
}

// CreateFromUser creates a new deck from the given user.
func (r *SQLRepository) CreateFromUser(ctx context.Context, user domain.User, deck *domain.Deck) error {
	return r.DBConn.WithContext(ctx).Model(&user).Association("OwnDecks").Append(deck)
}

// GetByID returns the deck with the given id.
func (r *SQLRepository) GetByID(ctx context.Context, id uint) (domain.Deck, error) {
	var deck domain.Deck
	err := r.DBConn.WithContext(ctx).Preload("Learners").First(&deck, id).Error
	return deck, err
}

// GetByUser returns the decks of the given user.
func (r *SQLRepository) GetByUser(ctx context.Context, user domain.User) ([]domain.Deck, error) {
	var decks []domain.Deck
	err := r.DBConn.WithContext(ctx).Model(&user).Association("OwnDecks").Find(&decks)
	return decks, err
}

// GetByLearner returns the decks of the given learner.
func (r *SQLRepository) GetByLearner(ctx context.Context, user domain.User) ([]domain.Deck, error) {
	var decks []domain.Deck
	err := r.DBConn.WithContext(ctx).Model(&user).Association("Learning").Find(&decks)
	return decks, err
}

// GetPublic returns the public decks.
func (r *SQLRepository) GetPublic(ctx context.Context) ([]domain.Deck, error) {
	var decks []domain.Deck
	err := r.DBConn.WithContext(ctx).Where("status = ?", domain.DeckStatusPublic).Find(&decks).Error
	return decks, err
}

// NewRepository returns a new repository.
func NewRepository(dbConn *gorm.DB) IRepository {
	return &SQLRepository{DBConn: dbConn}
}
