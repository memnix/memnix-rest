package deck

import (
	"context"

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
func (u *UseCase) GetByID(ctx context.Context, id uint) (domain.Deck, error) {
	var deckObject domain.Deck

	if cacheHit, _ := u.IRedisRepository.GetByID(ctx, id); cacheHit != "" {
		if err := cbor.Unmarshal([]byte(cacheHit), &deckObject); err == nil {
			return deckObject, nil
		}
	}

	deckObject, err := u.IRepository.GetByID(ctx, id)
	if err != nil {
		return domain.Deck{}, err
	}

	if marshalledDeck, err := cbor.Marshal(deckObject); err == nil {
		_ = u.IRedisRepository.SetByID(ctx, id, string(marshalledDeck))
	}

	return deckObject, nil
}

// Create creates a new deck.
func (u *UseCase) Create(ctx context.Context, deck *domain.Deck) error {
	if marshalledDeck, err := cbor.Marshal(deck); err == nil {
		_ = u.IRedisRepository.SetByID(ctx, deck.ID, string(marshalledDeck))
	}
	return u.IRepository.Create(ctx, deck)
}

// CreateFromUser creates a new deck from the given user.
func (u *UseCase) CreateFromUser(ctx context.Context, user domain.User, deck *domain.Deck) error {
	if marshalledDeck, err := cbor.Marshal(deck); err == nil {
		_ = u.IRedisRepository.SetByID(ctx, deck.ID, string(marshalledDeck))
	}
	return u.IRepository.CreateFromUser(ctx, user, deck)
}

// GetByUser returns the decks of the given user.
func (u *UseCase) GetByUser(ctx context.Context, user domain.User) ([]domain.Deck, error) {
	var ownedDecks []domain.Deck

	if cacheHit, _ := u.IRedisRepository.GetOwnedByUser(ctx, user.ID); cacheHit != "" {
		if err := cbor.Unmarshal([]byte(cacheHit), &ownedDecks); err == nil {
			return ownedDecks, nil
		}
	}

	ownedDecks, err := u.IRepository.GetByUser(ctx, user)
	if err != nil {
		return []domain.Deck{}, err
	}

	if marshalledDeck, err := cbor.Marshal(ownedDecks); err == nil {
		_ = u.IRedisRepository.SetOwnedByUser(ctx, user.ID, string(marshalledDeck))
	}
	return ownedDecks, nil
}

// GetByLearner returns the decks of the given learner.
func (u *UseCase) GetByLearner(ctx context.Context, user domain.User) ([]domain.Deck, error) {
	var learningDecks []domain.Deck

	if cacheHit, _ := u.IRedisRepository.GetLearningByUser(ctx, user.ID); cacheHit != "" {
		if err := cbor.Unmarshal([]byte(cacheHit), &learningDecks); err == nil {
			return learningDecks, nil
		}
	}

	learnedDecks, err := u.IRepository.GetByLearner(ctx, user)
	if err != nil {
		return []domain.Deck{}, err
	}

	if marshalledDeck, err := cbor.Marshal(learnedDecks); err == nil {
		_ = u.IRedisRepository.SetLearningByUser(ctx, user.ID, string(marshalledDeck))
	}
	return learnedDecks, nil
}

// GetPublic returns the public decks.
func (u *UseCase) GetPublic(ctx context.Context) ([]domain.Deck, error) {
	return u.IRepository.GetPublic(ctx)
}
