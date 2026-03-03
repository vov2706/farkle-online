package routes

import (
	"app/http/middlewares"
	"app/models"
	"app/services"
	"app/ws"
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

type PlayerToggleInput struct {
	Type    string          `json:"type"`
	IsReady json.RawMessage `json:"is_ready"`
}

func SetupBroadcasts(app *fiber.App, auth *services.AuthService) {
	hub := ws.NewHub()

	app.Use("/ws", middlewares.WebsocketProtocol(), middlewares.WsProtected(), func(c fiber.Ctx) error {
		user, err := auth.GetAuthUser(c)

		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals("current_user", user)

		return c.Next()
	})

	app.Get("/ws/lobby/:code", websocket.New(func(c *websocket.Conn) {
		code := c.Params("code")
		authUser := c.Locals("current_user").(*models.User)

		hub.Join(code, c)
		defer hub.Leave(code, c)

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			var in PlayerToggleInput
			if err = json.Unmarshal(msg, &in); err != nil {
				log.Println("fail to unmarshal: ", err)
				continue
			}

			out := map[string]interface{}{
				"type": in.Type,
				"data": map[string]any{
					"player_id": authUser.ID,
					"is_ready":  in.IsReady,
				},
			}
			body, _ := json.Marshal(out)

			hub.Broadcast(code, body)
		}
	}))
}
