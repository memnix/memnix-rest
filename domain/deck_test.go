package domain_test

import (
	"testing"

	"github.com/memnix/memnix-rest/domain"
)

func TestDeck_TableName(t *testing.T) {
	deck := &domain.Deck{}

	expected := "decks"
	tableName := deck.TableName()

	if tableName != expected {
		t.Errorf("Expected table name to be %s, but got %s", expected, tableName)
	}
}

func TestDeck_ToPublicDeck(t *testing.T) {
	deck := &domain.Deck{
		Name:        "Test Deck",
		Description: "Test Description",
		Lang:        "English",
		Banner:      "Test Banner",
	}

	deck.ID = 1

	expected := domain.PublicDeck{
		ID:          1,
		Name:        "Test Deck",
		Description: "Test Description",
		Lang:        "English",
		Banner:      "Test Banner",
	}

	publicDeck := deck.ToPublicDeck()

	if publicDeck != expected {
		t.Errorf("Expected public deck to be %v, but got %v", expected, publicDeck)
	}
}

func TestDeck_IsOwner(t *testing.T) {
	deck := &domain.Deck{
		OwnerID: 1,
	}

	// Test case 1: Owner ID matches
	id := uint(1)
	expected := true
	result := deck.IsOwner(id)
	if result != expected {
		t.Errorf("Expected IsOwner(%d) to be %v, but got %v", id, expected, result)
	}

	// Test case 2: Owner ID does not match
	id = uint(2)
	expected = false
	result = deck.IsOwner(id)
	if result != expected {
		t.Errorf("Expected IsOwner(%d) to be %v, but got %v", id, expected, result)
	}
}

func TestCreateDeck_ToDeck(t *testing.T) {
	createDeck := &domain.CreateDeck{
		Name:        "Test Name",
		Description: "Test Description",
		Lang:        "English",
	}

	expected := domain.Deck{
		Name:        "Test Name",
		Description: "Test Description",
		Lang:        "English",
		Status:      domain.DeckStatusPrivate,
	}

	deck := createDeck.ToDeck()

	if deck.Name != expected.Name {
		t.Errorf("Expected deck name to be %s, but got %s", expected.Name, deck.Name)
	}

	if deck.Description != expected.Description {
		t.Errorf("Expected deck description to be %s, but got %s", expected.Description, deck.Description)
	}

	if deck.Lang != expected.Lang {
		t.Errorf("Expected deck language to be %s, but got %s", expected.Lang, deck.Lang)
	}

	if deck.Status != expected.Status {
		t.Errorf("Expected deck status to be %s, but got %s", expected.Status, deck.Status)
	}
}

func TestDeckStatus_String(t *testing.T) {
	tests := []struct {
		expected string
		status   domain.DeckStatus
	}{
		{status: domain.DeckStatusPrivate, expected: "private"},
		{status: domain.DeckStatusToReview, expected: "to review"},
		{status: domain.DeckStatusPublic, expected: "public"},
	}

	for _, test := range tests {
		result := test.status.String()
		if result != test.expected {
			t.Errorf("Expected DeckStatus.String() to return %s, but got %s", test.expected, result)
		}
	}
}

func TestCreateDeck_Validate(t *testing.T) {
	testCases := []struct {
		createDeck *domain.CreateDeck
		name       string
		shouldErr  bool
	}{
		{
			name: "Valid CreateDeck",
			createDeck: &domain.CreateDeck{
				Name:        "Test Name",
				Description: "Test Description",
				Lang:        "EN",
				Banner:      "https://test.com/banner.png",
			},
			shouldErr: false,
		},
		{
			name: "No Name CreateDeck",
			createDeck: &domain.CreateDeck{
				Name:        "",
				Description: "Test Description",
				Lang:        "EN",
				Banner:      "https://test.com/banner.png",
			},
			shouldErr: true,
		},
		{
			name: "No Description CreateDeck",
			createDeck: &domain.CreateDeck{
				Name:        "Test Name",
				Description: "",
				Lang:        "EN",
				Banner:      "https://test.com/banner.png",
			},
			shouldErr: true,
		},
		{
			name: "No Lang CreateDeck",
			createDeck: &domain.CreateDeck{
				Name:        "Test Name",
				Description: "Test Description",
				Lang:        "",
				Banner:      "https://test.com/banner.png",
			},
			shouldErr: true,
		},
		{
			name: "No Banner CreateDeck",
			createDeck: &domain.CreateDeck{
				Name:        "Test Name",
				Description: "Test Description",
				Lang:        "EN",
				Banner:      "",
			},
			shouldErr: true,
		},
		{
			name: "Invalid Lang CreateDeck",
			createDeck: &domain.CreateDeck{
				Name:        "Test Name",
				Description: "Test Description",
				Lang:        "ENGLISH",
				Banner:      "https://test.com/banner.png",
			},
			shouldErr: true,
		},
		{
			name: "Invalid Banner CreateDeck",
			createDeck: &domain.CreateDeck{
				Name:        "Test Name",
				Description: "Test Description",
				Lang:        "EN",
				Banner:      "test.com/banner.png",
			},
			shouldErr: true,
		},
		{
			name: "Invalid Language CreateDeck",
			createDeck: &domain.CreateDeck{
				Name:        "Test Name",
				Description: "Test Description",
				Lang:        "INVALID",
				Banner:      "https://test.com/banner.png",
			},
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.createDeck.Validate()
			if tc.shouldErr && err == nil {
				t.Errorf("Expected error, but got nil")
			}
			if !tc.shouldErr && err != nil {
				t.Errorf("Expected no error, but got %v", err)
			}
		})
	}
}
