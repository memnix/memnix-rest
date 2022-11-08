package viewmodels

import (
	"github.com/memnix/memnixrest/models"
)

// DeckResponse structure
type DeckResponse struct {
	DeckID uint           `json:"deck_id"`
	Deck   models.Deck    `json:"deck"`
	Cards  []ResponseCard `json:"cards"`
	Count  int            `json:"count"`
}

type TodayResponse struct {
	DecksReponses []DeckResponse `json:"decks_responses"`
	Count         int            `json:"count"`
}

func (today *TodayResponse) AppendDeckResponse(deckResponse DeckResponse) {
	today.DecksReponses = append(today.DecksReponses, deckResponse)
}
