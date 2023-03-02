package deck

import "github.com/memnix/memnix-rest/domain"

type UseCase struct {
	IRepository
}

func NewUseCase(repo IRepository) IUseCase {
	return &UseCase{IRepository: repo}
}

// GetByID returns the deck with the given id.
func (u *UseCase) GetByID(id uint) (domain.Deck, error) {
	return u.IRepository.GetByID(id)
}

func (u *UseCase) Create(deck *domain.Deck) error {
	return u.IRepository.Create(deck)
}

func (u *UseCase) CreateFromUser(user domain.User, deck *domain.Deck) error {
	return u.IRepository.CreateFromUser(user, deck)
}

// GetByUser returns the decks of the given user.
func (u *UseCase) GetByUser(user domain.User) ([]domain.Deck, error) {
	return u.IRepository.GetByUser(user)
}

// GetByLearner returns the decks of the given learner.
func (u *UseCase) GetByLearner(user domain.User) ([]domain.Deck, error) {
	return u.IRepository.GetByLearner(user)
}

// GetPublic returns the public decks.
func (u *UseCase) GetPublic() ([]domain.Deck, error) {
	return u.IRepository.GetPublic()
}
