package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/internal/deck"
	"github.com/memnix/memnix-rest/pkg/utils"
	"github.com/memnix/memnix-rest/views"
	"github.com/pkg/errors"
)

// DeckController is the controller for the deck routes
type DeckController struct {
	deck.IUseCase
}

// NewDeckController returns a new deck controller
func NewDeckController(useCase deck.IUseCase) DeckController {
	return DeckController{IUseCase: useCase}
}

// GetByID is the controller for the get deck by id route
//
//	@Summary		Get deck by id
//	@Description	Get deck by id
//	@Tags			Deck
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint	true	"Deck ID"
//	@Success		200	{object}	views.HTTPResponseVM
//	@Failure		400	{object}	views.HTTPResponseVM
//	@Failure		403	{object}	views.HTTPResponseVM
//	@Failure		404	{object}	views.HTTPResponseVM
//	@Router			/v2/deck/{id} [get]
//	@Security		Bearer
func (d *DeckController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	uintID, err := utils.ConvertStrToUInt(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(views.NewHTTPResponseVMFromError(err))
	}
	deckObject, err := d.IUseCase.GetByID(c.UserContext(), uintID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(views.NewHTTPResponseVMFromError(err))
	}

	if !deckObject.IsOwner(utils.GetUserFromContext(c).ID) && utils.GetUserFromContext(c).Permission != domain.PermissionAdmin {
		return c.Status(fiber.StatusForbidden).JSON(views.NewHTTPResponseVMFromError(errors.New("deck is private")))
	}

	return c.Status(fiber.StatusOK).JSON(views.NewHTTPResponseVM("deck found", deckObject))
}

// Create is the controller for the create deck route
//
//	@Summary		Create deck
//	@Description	Create deck
//	@Tags			Deck
//	@Accept			json
//	@Produce		json
//	@Param			deck	body		domain.Deck	true	"Deck object"
//	@Success		201		{object}	views.HTTPResponseVM
//	@Failure		400		{object}	views.HTTPResponseVM
//	@Router			/v2/deck [post]
//	@Security		Bearer
func (d *DeckController) Create(c *fiber.Ctx) error {
	var createDeck domain.CreateDeck
	err := c.BodyParser(&createDeck)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(views.NewHTTPResponseVMFromError(err))
	}

	if err := createDeck.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(views.NewHTTPResponseVMFromError(err))
	}

	deckObject := createDeck.ToDeck()

	if err := d.IUseCase.CreateFromUser(c.UserContext(), *utils.GetUserFromContext(c), &deckObject); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(views.NewHTTPResponseVMFromError(err))
	}

	return c.Status(fiber.StatusCreated).JSON(views.NewHTTPResponseVM("deck created", deckObject.ToPublicDeck()))
}

// GetOwned is the controller for the get owned decks route
//
//	@Summary		Get owned decks
//	@Description	Get owned decks
//	@Tags			Deck
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	views.HTTPResponseVM
//	@Failure		400	{object}	views.HTTPResponseVM
//	@Failure		404	{object}	views.HTTPResponseVM
//	@Router			/v2/deck/owned [get]
//	@Security		Bearer
func (d *DeckController) GetOwned(c *fiber.Ctx) error {
	deckObjects, err := d.IUseCase.GetByUser(c.UserContext(), *utils.GetUserFromContext(c))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(views.NewHTTPResponseVMFromError(err))
	}

	return c.Status(fiber.StatusOK).JSON(views.NewHTTPResponseVM("owned decks found", deckObjects))
}

// GetLearning is the controller for the get learning decks route
//
//	@Summary		Get learning decks
//	@Description	Get learning decks
//	@Tags			Deck
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	views.HTTPResponseVM
//	@Failure		400	{object}	views.HTTPResponseVM
//	@Failure		404	{object}	views.HTTPResponseVM
//	@Router			/v2/deck/learning [get]
//	@Security		Bearer
func (d *DeckController) GetLearning(c *fiber.Ctx) error {
	deckObjects, err := d.IUseCase.GetByLearner(c.UserContext(), *utils.GetUserFromContext(c))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(views.NewHTTPResponseVMFromError(err))
	}

	return c.Status(fiber.StatusOK).JSON(views.NewHTTPResponseVM("learning decks found", deck.ConvertToPublic(deckObjects)))
}

// GetPublic is the controller for the get public decks route
//
//	@Summary		Get public decks
//	@Description	Get public decks
//	@Tags			Deck
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	views.HTTPResponseVM
//	@Failure		400	{object}	views.HTTPResponseVM
//	@Failure		404	{object}	views.HTTPResponseVM
//	@Router			/v2/deck/public [get]
//	@Security		Bearer
func (d *DeckController) GetPublic(c *fiber.Ctx) error {
	deckObjects, err := d.IUseCase.GetPublic(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(views.NewHTTPResponseVMFromError(err))
	}

	return c.Status(fiber.StatusOK).JSON(views.NewHTTPResponseVM("public decks found", deckObjects))
}
