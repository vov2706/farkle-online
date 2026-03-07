package ws

import (
	"app/services"
	"app/ws"
	"encoding/json"
)

func StartGame(service *services.GameService) ws.Handler {
	return func(ctx *ws.Ctx) error {
		var data map[string]any
		if err := json.Unmarshal(ctx.Env.Data, &data); err != nil {
			return nil
		}

		code := data["code"].(string)
		game, err := service.StartGame(ctx.User, code)

		if err != nil {
			return err
		}

		ctx.Hub.Broadcast(ctx.Env.Channel, ws.Event("game_started", ctx.Env.Channel, map[string]any{
			"start_time": game.StartedAt.UnixMilli(),
		}))

		return nil
	}
}
