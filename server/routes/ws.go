package routes

import (
	"app/container"
	"app/http/middlewares"
	"app/models"
	"app/ws"
	ws2 "app/ws/handlers"
	"log"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func SetupBroadcasts(app *fiber.App, ctn *container.Container) {
	ctn.WS.Authorizer.Channel("presence-lobby:{code}", []string{"code"}, func(user any, params map[string]string) (bool, map[string]any) {
		u := user.(*models.User)
		code := params["code"]
		_ = code

		return true, fiber.Map{
			"username": u.Username,
		}
	})

	ctn.WS.Router.On("subscribe", ws2.SubscribeHandler(ctn.WS.Authorizer))
	ctn.WS.Router.On("unsubscribe", ws2.UnsubscribeHandler())
	ctn.WS.Router.On("lobby.player_ready", ws2.LobbyPlayerReadyHandler())
	ctn.WS.Router.On("start_game", ws2.StartGame(ctn.Services.Game))
	ctn.WS.Router.On("lobby.player_left", ws2.LobbyPlayerLeft(ctn.Repos.Game))

	app.Use("/ws", middlewares.WebsocketProtocol(), middlewares.WsProtected(), func(c fiber.Ctx) error {
		user, err := ctn.Services.Auth.GetAuthUser(c)

		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals("current_user", user)

		return c.Next()
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		user := c.Locals("current_user").(*models.User)

		defer func() {
			_, presenceLeft := ctn.WS.Hub.LeaveAll(c)
			for ch, u := range presenceLeft {
				ctn.WS.Hub.Broadcast(ch, ws.Event("presence.leaving", ch, map[string]any{
					"user": u,
				}))
			}
		}()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			ctn.WS.Router.Handle(c, user, msg)
		}
	}))
}
