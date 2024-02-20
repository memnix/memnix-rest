package controllers

import "github.com/memnix/memnix-rest/services/mcq"

type McqController struct {
	mcq.IUseCase
}

func NewMcqController(useCase mcq.IUseCase) McqController {
	return McqController{useCase}
}
