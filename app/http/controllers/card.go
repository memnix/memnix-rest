package controllers

import "github.com/memnix/memnix-rest/internal/card"

type CardController struct {
	card.IUseCase
}

func NewCardController(useCase card.IUseCase) CardController {
	return CardController{IUseCase: useCase}
}
