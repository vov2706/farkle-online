package main

import (
	"app/config"
	"app/database"
	"app/repositories"
	"app/routes"
	"app/services"
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
	routes.SetupBroadcasts(app, services.NewAuthService(repositories.NewUserRepository(database.DB)))
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":" + config.Config("API_PORT")))
}
