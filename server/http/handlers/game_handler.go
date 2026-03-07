package handlers

import (
	"app/http/inputs"
	"app/http/responses"
	"app/services"
	"app/utils"
	"app/ws"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type GameHandler struct {
	gameService *services.GameService
	auth        *services.AuthService
	hub         *ws.Hub
}

func NewGameHandler(
	gameService *services.GameService,
	authService *services.AuthService,
	hub *ws.Hub,
) *GameHandler {
	return &GameHandler{
		gameService: gameService,
		auth:        authService,
		hub:         hub,
	}
}

func (handler *GameHandler) Store(c fiber.Ctx) error {
	input := new(inputs.CreateGameInput)

	if err := c.Bind().Body(input); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"field":   e.Field(),
					"message": e.Error(),
				})
			}
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// custom validation
	if err := input.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	authUser, err := handler.auth.GetAuthUser(c)

	if err != nil {
		return err
	}

	game, err := handler.gameService.CreateGame(authUser, input)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(responses.NewGameResource(*game))
}

func (handler *GameHandler) Index(c fiber.Ctx) error {
	page, err := utils.ConvertQueryParamToUint(c, "page")

	if err != nil {
		return err
	}

	authUser, err := handler.auth.GetAuthUser(c)

	if err != nil {
		return err
	}

	results, err := handler.gameService.PaginatePublicGames(
		authUser.ID,
		page,
		10,
		c.Query("search"),
	)

	if err != nil {
		return err
	}

	return c.JSON(results)
}

func (handler *GameHandler) Show(c fiber.Ctx) error {
	code := c.Params("code")
	game := handler.gameService.GetGameByCode(code)

	if game == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "game not found",
		})
	}

	return c.JSON(responses.NewGameResource(*game))
}

func (handler *GameHandler) Leave(c fiber.Ctx) error {
	authUser, err := handler.auth.GetAuthUser(c)

	if err != nil {
		return err
	}

	err = handler.gameService.LeaveCurrentGame(authUser)

	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}

func (handler *GameHandler) GetCurrentGame(c fiber.Ctx) error {
	authUser, err := handler.auth.GetAuthUser(c)

	if err != nil {
		return err
	}

	game := handler.gameService.GetCurrentGame(authUser)

	if game != nil {
		return c.JSON(responses.NewGameResource(*game))
	}

	return c.JSON(nil)
}

func (handler *GameHandler) Join(c fiber.Ctx) error {
	code := c.Params("code")
	authUser, err := handler.auth.GetAuthUser(c)

	if err != nil {
		return err
	}

	game, err := handler.gameService.JoinToGame(authUser, code)

	if err != nil {
		return err
	}

	ch := "presence-lobby:" + game.Code
	handler.hub.Broadcast(ch, ws.Event("lobby.player_joined", ch, fiber.Map{
		"player": authUser.Username,
	}))

	return c.JSON(fiber.Map{
		"game": responses.NewGameResource(*game),
	})
}
