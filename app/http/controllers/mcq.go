package controllers

import "github.com/memnix/memnix-rest/internal/mcq"

type McqController struct {
	mcq.IUseCase
}

func NewMcqController(useCase mcq.IUseCase) McqController {
	return McqController{useCase}
}
