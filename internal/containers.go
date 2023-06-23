package internal

import (
	"sync"

	"github.com/memnix/memnix-rest/app/http/controllers"
)

// ServiceContainer is the service container interface.
type ServiceContainer interface {
	GetUser() controllers.UserController   // GetUser returns the user controller.
	GetAuth() controllers.AuthController   // GetAuth returns the auth controller.
	GetJwt() controllers.JwtController     // GetJwt returns the jwt controller.
	GetOAuth() controllers.OAuthController // GetOAuth returns the oauth controller.
	GetDeck() controllers.DeckController   // GetDeck returns the deck controller.
}

type kernel struct{}

// GetUser returns the user controller.
func (kernel) GetUser() controllers.UserController {
	return InitializeUser()
}

// GetAuth returns the auth controller.
func (kernel) GetAuth() controllers.AuthController {
	return InitializeAuth()
}

// GetJwt returns the jwt controller.
func (kernel) GetJwt() controllers.JwtController {
	return InitializeJWT()
}

// GetOAuth returns the oauth controller.
func (kernel) GetOAuth() controllers.OAuthController {
	return InitializeOAuth()
}

// GetDeck returns the deck controller.
func (kernel) GetDeck() controllers.DeckController {
	return InitializeDeck()
}

func (kernel) GetCard() controllers.CardController {
	return InitializeCard()
}

func (kernel) GetMcq() controllers.McqController {
	return InitializeMcq()
}

var (
	k             *kernel   // kernel is the service container
	containerOnce sync.Once // containerOnce is used to ensure that the service container is only initialized once
)

// GetServiceContainer returns the service container
func GetServiceContainer() ServiceContainer {
	containerOnce.Do(func() {
		k = &kernel{}
	})
	return k
}
