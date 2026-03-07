package handlers

import (
	"app/http/responses"
	"app/services"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	auth        *services.AuthService
	gameService *services.GameService
}

func NewUserHandler(auth *services.AuthService, gameService *services.GameService) *UserHandler {
	return &UserHandler{auth: auth, gameService: gameService}
}

func (handler *UserHandler) GetProfile(c fiber.Ctx) error {
	user, err := handler.auth.GetAuthUser(c, "Balances.Currency")

	if err != nil {
		return err
	}

	user.CurrentGame = handler.gameService.GetCurrentGame(user)

	return c.JSON(responses.UserResponse{
		Data: responses.NewUserResource(*user),
	})
}
