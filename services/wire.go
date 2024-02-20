//go:build wireinject
// +build wireinject

package services

import (
	"github.com/google/wire"
	"github.com/memnix/memnix-rest/app/v1/controllers"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/services/auth"
	"github.com/memnix/memnix-rest/services/card"
	"github.com/memnix/memnix-rest/services/deck"
	"github.com/memnix/memnix-rest/services/mcq"
	"github.com/memnix/memnix-rest/services/user"
)

// InitializeUser initializes the user controller.
func InitializeUser() controllers.UserController {
	wire.Build(controllers.NewUserController, user.NewUseCase, user.NewRedisRepository, infrastructures.GetRedisClient, user.NewRepository, infrastructures.GetDBConn, user.NewRistrettoCache, infrastructures.GetRistrettoCache)
	return controllers.UserController{}
}

// InitializeAuth initializes the auth controller.
func InitializeAuth() controllers.AuthController {
	wire.Build(controllers.NewAuthController, auth.NewUseCase, user.NewRepository, infrastructures.GetDBConn)
	return controllers.AuthController{}
}

// InitializeJWT initializes the jwt controller.
func InitializeJWT() controllers.JwtController {
	wire.Build(controllers.NewJwtController, user.NewUseCase, user.NewRedisRepository, infrastructures.GetRedisClient, user.NewRepository, infrastructures.GetDBConn, user.NewRistrettoCache, infrastructures.GetRistrettoCache)
	return controllers.JwtController{}
}

// InitializeOAuth initializes the oauth controller.
func InitializeOAuth() controllers.OAuthController {
	wire.Build(controllers.NewOAuthController, auth.NewUseCase, user.NewRepository, infrastructures.GetDBConn, auth.NewRedisRepository, infrastructures.GetRedisClient)
	return controllers.OAuthController{}
}

// InitializeDeck initializes the deck controller.
func InitializeDeck() controllers.DeckController {
	wire.Build(controllers.NewDeckController, deck.NewUseCase, deck.NewRepository, infrastructures.GetDBConn, deck.NewRedisRepository, infrastructures.GetRedisClient)
	return controllers.DeckController{}
}

// InitializeCard initializes the card controller.
func InitializeCard() controllers.CardController {
	wire.Build(controllers.NewCardController, card.NewUseCase, card.NewRepository, infrastructures.GetDBConn, card.NewRedisRepository, infrastructures.GetRedisClient)
	return controllers.CardController{}
}

// InitializeMcq initializes the mcq controller.
func InitializeMcq() controllers.McqController {
	wire.Build(controllers.NewMcqController, mcq.NewUseCase, mcq.NewRepository, infrastructures.GetDBConn, mcq.NewRedisRepository, infrastructures.GetRedisClient)
	return controllers.McqController{}
}
