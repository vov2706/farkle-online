package ws

import (
	"app/repositories"
	"app/ws"
	"encoding/json"

	"github.com/gofiber/fiber/v3"
)

type PlayerReadyInput struct {
	IsReady bool `json:"is_ready"`
}

func LobbyPlayerReadyHandler() ws.Handler {
	return func(ctx *ws.Ctx) error {
		if ctx.Env.Channel == "" {
			return nil
		}

		var in PlayerReadyInput
		if err := json.Unmarshal(ctx.Env.Data, &in); err != nil {
			return nil
		}

		// todo: check if user exists in game

		ctx.Hub.Broadcast(ctx.Env.Channel, ws.Event("lobby.player_ready", ctx.Env.Channel, map[string]any{
			"player":   ctx.User.Username,
			"is_ready": in.IsReady,
		}))
		return nil
	}
}

func LobbyPlayerLeft(repo *repositories.GameRepository) ws.Handler {
	return func(ctx *ws.Ctx) error {
		var data map[string]any
		if err := json.Unmarshal(ctx.Env.Data, &data); err != nil {
			return nil
		}

		code := data["code"].(string)
		game, err := repo.FindGameByCode(code)

		if err != nil {
			return nil
		}

		if ctx.User.ID == game.CreatorID {
			err = repo.DeleteGame(ctx.User, game)

			ctx.Hub.Broadcast(ctx.Env.Channel, ws.Event("lobby.game_deleted", ctx.Env.Channel, fiber.Map{
				"code": code,
			}))

			return nil
		}

		err = repo.DetachUserFromGame(ctx.User, game)
		ctx.Hub.Broadcast(ctx.Env.Channel, ws.Event("lobby.player_left", ctx.Env.Channel, fiber.Map{
			"player": ctx.User.Username,
		}))

		return nil
	}
}
