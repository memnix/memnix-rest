package internal

import (
	"sync"

	"github.com/memnix/memnix-rest/app/http/controllers"
)

type ServiceContainer struct {
	user  controllers.UserController
	auth  controllers.AuthController
	jwt   controllers.JwtController
	oAuth controllers.OAuthController
	deck  controllers.DeckController
	card  controllers.CardController
	mcq   controllers.McqController
}

var (
	container *ServiceContainer //nolint:gochecknoglobals // Singleton
	once      sync.Once         //nolint:gochecknoglobals // Singleton
)

func DefaultServiceContainer() *ServiceContainer {
	userController := InitializeUser()
	authController := InitializeAuth()
	jwtController := InitializeJWT()
	oAuthController := InitializeOAuth()
	deckController := InitializeDeck()
	cardController := InitializeCard()
	mcqController := InitializeMcq()

	return NewServiceContainer(userController, authController,
		jwtController, oAuthController, deckController,
		cardController, mcqController)
}

func NewServiceContainer(user controllers.UserController, auth controllers.AuthController, jwt controllers.JwtController, oAuth controllers.OAuthController, deck controllers.DeckController, card controllers.CardController, mcq controllers.McqController) *ServiceContainer {
	once.Do(func() {
		container = &ServiceContainer{
			user:  user,
			auth:  auth,
			jwt:   jwt,
			oAuth: oAuth,
			deck:  deck,
			card:  card,
			mcq:   mcq,
		}
	})
	return container
}

func (sc *ServiceContainer) User() controllers.UserController {
	return sc.user
}

func (sc *ServiceContainer) Auth() controllers.AuthController {
	return sc.auth
}

func (sc *ServiceContainer) Jwt() controllers.JwtController {
	return sc.jwt
}

func (sc *ServiceContainer) OAuth() controllers.OAuthController {
	return sc.oAuth
}

func (sc *ServiceContainer) Deck() controllers.DeckController {
	return sc.deck
}

func (sc *ServiceContainer) Card() controllers.CardController {
	return sc.card
}

func (sc *ServiceContainer) Mcq() controllers.McqController {
	return sc.mcq
}
