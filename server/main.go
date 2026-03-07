package main

import (
	"app/config"
	"app/container"
	"app/database"
	"app/routes"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		ServerHeader:  "Fiber",
		AppName:       "Tavern Dice",
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowCredentials: false,
	}))

	database.Connect()

	ctn := container.New()

	routes.SetupBroadcasts(app, ctn)
	routes.SetupRoutes(app, ctn)

	log.Fatal(app.Listen(":" + config.Config("API_PORT")))
}
