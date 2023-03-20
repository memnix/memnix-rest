package deck

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/memnix/memnix-rest/domain"
)

// UseCase is the deck use case.
type UseCase struct {
	IRepository
	IRedisRepository
}

// NewUseCase returns a new deck use case.
func NewUseCase(repo IRepository, redis IRedisRepository) IUseCase {
	return &UseCase{IRepository: repo, IRedisRepository: redis}
}

// GetByID returns the deck with the given id.
func (u *UseCase) GetByID(id uint) (domain.Deck, error) {
	var deckObject domain.Deck

	if cacheHit, _ := u.IRedisRepository.GetByID(id); cacheHit != "" {
		if err := cbor.Unmarshal([]byte(cacheHit), &deckObject); err == nil {
			return deckObject, nil
		}
	}

	deckObject, err := u.IRepository.GetByID(id)
	if err != nil {
		return domain.Deck{}, err
	}

	if marshalledDeck, err := cbor.Marshal(deckObject); err == nil {
		_ = u.IRedisRepository.SetByID(id, string(marshalledDeck))
	}

	return deckObject, nil
}

// Create creates a new deck.
func (u *UseCase) Create(deck *domain.Deck) error {
	if marshalledDeck, err := cbor.Marshal(deck); err == nil {
		_ = u.IRedisRepository.SetByID(deck.ID, string(marshalledDeck))
	}
	return u.IRepository.Create(deck)
}

// CreateFromUser creates a new deck from the given user.
func (u *UseCase) CreateFromUser(user domain.User, deck *domain.Deck) error {
	if marshalledDeck, err := cbor.Marshal(deck); err == nil {
		_ = u.IRedisRepository.SetByID(deck.ID, string(marshalledDeck))
	}
	return u.IRepository.CreateFromUser(user, deck)
}

// GetByUser returns the decks of the given user.
func (u *UseCase) GetByUser(user domain.User) ([]domain.Deck, error) {
	var ownedDecks []domain.Deck

	if cacheHit, _ := u.IRedisRepository.GetOwnedByUser(user.ID); cacheHit != "" {
		if err := cbor.Unmarshal([]byte(cacheHit), &ownedDecks); err == nil {
			return ownedDecks, nil
		}
	}

	ownedDecks, err := u.IRepository.GetByUser(user)
	if err != nil {
		return []domain.Deck{}, err
	}

	if marshalledDeck, err := cbor.Marshal(ownedDecks); err == nil {
		_ = u.IRedisRepository.SetOwnedByUser(user.ID, string(marshalledDeck))
	}
	return ownedDecks, nil
}

// GetByLearner returns the decks of the given learner.
func (u *UseCase) GetByLearner(user domain.User) ([]domain.Deck, error) {
	var learningDecks []domain.Deck

	if cacheHit, _ := u.IRedisRepository.GetLearningByUser(user.ID); cacheHit != "" {
		if err := cbor.Unmarshal([]byte(cacheHit), &learningDecks); err == nil {
			return learningDecks, nil
		}
	}

	learnedDecks, err := u.IRepository.GetByLearner(user)
	if err != nil {
		return []domain.Deck{}, err
	}

	if marshalledDeck, err := cbor.Marshal(learnedDecks); err == nil {
		_ = u.IRedisRepository.SetLearningByUser(user.ID, string(marshalledDeck))
	}
	return learnedDecks, nil
}

// GetPublic returns the public decks.
func (u *UseCase) GetPublic() ([]domain.Deck, error) {
	return u.IRepository.GetPublic()
}
