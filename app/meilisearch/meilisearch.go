package meilisearch

import (
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/internal/deck"
)

type MeiliSearch struct {
	deck.IUseCase
}

func NewMeiliSearch(deck deck.IUseCase) MeiliSearch {
	return MeiliSearch{
		IUseCase: deck,
	}
}

func (m *MeiliSearch) CreateSearchIndex() error {
	decks, err := m.GetPublic()
	if err != nil {
		return err
	}

	decksIndex := make([]domain.DeckIndex, len(decks))

	for idx, deckModel := range decks {
		decksIndex[idx] = domain.DeckIndex{
			"id":          deckModel.ID,
			"title":       deckModel.Name,
			"description": deckModel.Description,
			"lang":        deckModel.Lang,
		}
	}

	_, err = infrastructures.GetMeiliSearchClient().Index("decks").AddDocuments(decksIndex)
	if err != nil {
		return err
	}

	return nil
}

func InitMeiliSearch(meiliSearch MeiliSearch) error {
	err := meiliSearch.CreateSearchIndex()
	if err != nil {
		return err
	}

	return nil
}
