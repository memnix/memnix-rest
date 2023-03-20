package meilisearch

import (
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/internal/deck"
)

// MeiliSearch is a struct that implements the deck.IUseCase interface.
type MeiliSearch struct {
	deck.IUseCase
}

// NewMeiliSearch creates a new MeiliSearch instance.
func NewMeiliSearch(deck deck.IUseCase) MeiliSearch {
	return MeiliSearch{
		IUseCase: deck,
	}
}

// CreateSearchIndex creates a new index in MeiliSearch.
func (m *MeiliSearch) CreateSearchIndex() error {
	decks, err := m.GetPublic()
	if err != nil {
		return err
	}

	decksIndex := make([]domain.DeckIndex, len(decks))

	for idx, deckModel := range decks {
		decksIndex[idx] = domain.DeckIndex{
			"id":          deckModel.ID,
			"name":        deckModel.Name,
			"description": deckModel.Description,
			"lang":        deckModel.Lang,
			"banner":      deckModel.Banner,
		}
	}

	_, err = infrastructures.GetMeiliSearchClient().Index("decks").AddDocuments(decksIndex)
	if err != nil {
		return err
	}

	return nil
}

// InitMeiliSearch initializes the MeiliSearch instance.
func InitMeiliSearch(meiliSearch MeiliSearch) error {
	err := meiliSearch.CreateSearchIndex()
	if err != nil {
		return err
	}

	return nil
}
