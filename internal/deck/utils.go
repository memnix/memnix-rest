package deck

import "github.com/memnix/memnix-rest/domain"

// ConvertToPublic converts a slice of domain.Deck to a slice of domain.PublicDeck.
func ConvertToPublic(deck []domain.Deck) []domain.PublicDeck {
	decksPublic := make([]domain.PublicDeck, 0, len(deck))
	for idx, d := range deck {
		decksPublic[idx] = d.ToPublicDeck()
	}
	return decksPublic
}
