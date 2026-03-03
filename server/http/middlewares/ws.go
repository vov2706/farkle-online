package middlewares

import (
	"app/config"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
)

func WebsocketProtocol() fiber.Handler {
	return func(c fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}

func WsProtected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		Extractor:  extractors.FromQuery("token"),
		SigningKey: jwtware.SigningKey{Key: []byte(config.Config("JWT_SECRET"))},
	})
}
