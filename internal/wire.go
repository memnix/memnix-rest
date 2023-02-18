//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/memnix/memnix-rest/app/http/controllers"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/internal/barista"
	"github.com/memnix/memnix-rest/internal/kliento"
)

func InitializeBarista() controllers.BaristaController {
	wire.Build(controllers.NewBaristaController, barista.NewUseCase, barista.NewPgRepository, infrastructures.GetDBConn)
	return controllers.BaristaController{}
}

func InitializeKliento() controllers.KlientoController {
	wire.Build(controllers.NewKlientoController, kliento.NewUseCase, kliento.NewRedisRepository, infrastructures.GetRedisClient)
	return controllers.KlientoController{}
}
