package services

import (
	"github.com/memnix/memnixrest/app/controllers"
	"github.com/memnix/memnixrest/data/infrastructures"
	"github.com/memnix/memnixrest/data/repositories"
	"github.com/memnix/memnixrest/interfaces"
)

type McqService struct {
	interfaces.IMcqRepository
}

func (k *kernel) InjectMcqController() controllers.McqController {
	DBConn := infrastructures.GetDBConn()

	mcqRepository := &repositories.McqRepository{DBConn: DBConn}
	mcqService := &McqService{IMcqRepository: mcqRepository}
	mcqController := controllers.McqController{IMcqService: mcqService}

	return mcqController
}
