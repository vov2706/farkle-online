package routes

import (
	"app/container"
	"app/http/middlewares"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func SetupRoutes(app *fiber.App, ctn *container.Container) {
	// Middleware
	api := app.Group("/api", logger.New())

	// Auth
	api.Post("/login", ctn.Handlers.Auth.Login)
	api.Post("/register", ctn.Handlers.Auth.Register)

	// Profile
	api.Get("/profile", middlewares.Protected(), ctn.Handlers.User.GetProfile)

	// Games
	gameGroup := api.Group("/games", middlewares.Protected())
	gameGroup.Get("/", ctn.Handlers.Game.Index)
	gameGroup.Get("/current", ctn.Handlers.Game.GetCurrentGame)
	gameGroup.Get("/:code", ctn.Handlers.Game.Show)
	gameGroup.Post("/", ctn.Handlers.Game.Store)
	gameGroup.Delete("/leave", ctn.Handlers.Game.Leave)
	gameGroup.Post("/join/:code", ctn.Handlers.Game.Join)

	// Currencies
	api.Get("/currencies", ctn.Handlers.Currency.Index)

	dist := "../client/dist"

	// Frontend assets
	app.Use("/assets", static.New(filepath.Join(dist, "assets")))
	app.Use("/images", static.New(filepath.Join(dist, "images")))
	app.Use("/favicon.ico", static.New(filepath.Join(dist, "favicon.ico")))

	// Frontend
	app.Get("*", func(c fiber.Ctx) error {
		path := c.Path()

		if strings.HasPrefix(path, "/api") {
			return fiber.ErrNotFound
		}

		return c.SendFile(filepath.Join(dist, "index.html"))
	})

	// 404
	app.Use(func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound) // => 404 "Not Found"
	})
}
