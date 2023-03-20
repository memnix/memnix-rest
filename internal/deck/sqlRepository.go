package deck

import (
	"github.com/memnix/memnix-rest/domain"
	"gorm.io/gorm"
)

// SQLRepository is the repository for the deck.
type SQLRepository struct {
	DBConn *gorm.DB
}

// Create creates a new deck.
func (r *SQLRepository) Create(deck *domain.Deck) error {
	return r.DBConn.Create(&deck).Error
}

// Update updates the deck with the given id.
func (r *SQLRepository) Update(deck *domain.Deck) error {
	return r.DBConn.Save(&deck).Error
}

// Delete deletes the deck with the given id.
func (r *SQLRepository) Delete(id uint) error {
	return r.DBConn.Delete(&domain.Deck{}, id).Error
}

// CreateFromUser creates a new deck from the given user.
func (r *SQLRepository) CreateFromUser(user domain.User, deck *domain.Deck) error {
	return r.DBConn.Model(&user).Association("OwnDecks").Append(deck)
}

// GetByID returns the deck with the given id.
func (r *SQLRepository) GetByID(id uint) (domain.Deck, error) {
	var deck domain.Deck
	err := r.DBConn.Preload("Learners").First(&deck, id).Error
	return deck, err
}

// GetByUser returns the decks of the given user.
func (r *SQLRepository) GetByUser(user domain.User) ([]domain.Deck, error) {
	var decks []domain.Deck
	err := r.DBConn.Model(&user).Association("OwnDecks").Find(&decks)
	return decks, err
}

// GetByLearner returns the decks of the given learner.
func (r *SQLRepository) GetByLearner(user domain.User) ([]domain.Deck, error) {
	var decks []domain.Deck
	err := r.DBConn.Model(&user).Association("Learning").Find(&decks)
	return decks, err
}

// GetPublic returns the public decks.
func (r *SQLRepository) GetPublic() ([]domain.Deck, error) {
	var decks []domain.Deck
	err := r.DBConn.Where("status = ?", domain.DeckStatusPublic).Find(&decks).Error
	return decks, err
}

// NewRepository returns a new repository.
func NewRepository(dbConn *gorm.DB) IRepository {
	return &SQLRepository{DBConn: dbConn}
}
