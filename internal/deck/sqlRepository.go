package deck

import (
	"github.com/memnix/memnix-rest/domain"
	"gorm.io/gorm"
)

type SqlRepository struct {
	DBConn *gorm.DB
}

func (r *SqlRepository) Create(deck *domain.Deck) error {
	return r.DBConn.Create(&deck).Error
}

func (r *SqlRepository) Update(deck *domain.Deck) error {
	return r.DBConn.Save(&deck).Error
}

func (r *SqlRepository) Delete(id uint) error {
	return r.DBConn.Delete(&domain.Deck{}, id).Error
}

func (r *SqlRepository) CreateFromUser(user domain.User, deck *domain.Deck) error {
	return r.DBConn.Model(&user).Association("OwnDecks").Append(deck)
}

// GetByID returns the deck with the given id.
func (r *SqlRepository) GetByID(id uint) (domain.Deck, error) {
	var deck domain.Deck
	err := r.DBConn.Preload("Learners").First(&deck, id).Error
	return deck, err
}

func (r *SqlRepository) GetByUser(user domain.User) ([]domain.Deck, error) {
	var decks []domain.Deck
	err := r.DBConn.Model(&user).Association("OwnDecks").Find(&decks)
	return decks, err
}

func (r *SqlRepository) GetByLearner(user domain.User) ([]domain.Deck, error) {
	var decks []domain.Deck
	err := r.DBConn.Model(&user).Association("Learning").Find(&decks)
	return decks, err
}

func (r *SqlRepository) GetPublic() ([]domain.Deck, error) {
	var decks []domain.Deck
	err := r.DBConn.Where("status = ?", domain.DeckStatusPublic).Find(&decks).Error
	return decks, err
}

func NewRepository(dbConn *gorm.DB) IRepository {
	return &SqlRepository{DBConn: dbConn}
}
