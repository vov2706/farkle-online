package routes

import (
	"app/database"
	"app/http/handlers"
	"app/http/middlewares"
	"app/repositories"
	"app/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())

	// Init repositories
	balanceRepo := repositories.NewBalanceRepository(database.DB)
	gameRepo := repositories.NewGameRepository(database.DB)
	userRepo := repositories.NewUserRepository(database.DB)
	currencyRepo := repositories.NewCurrencyRepository(database.DB)

	// Init services
	gameService := services.NewGameService(balanceRepo, currencyRepo, gameRepo)
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo, currencyRepo, balanceRepo)

	// Init handlers
	gameHandler := handlers.NewGameHandler(gameService, authService)
	authHandler := handlers.NewAuthHandler(authService, userService)
	userHandler := handlers.NewUserHandler(authService, gameService)
	currencyHandler := handlers.NewCurrencyHandler(currencyRepo)

	// Auth
	api.Post("/login", authHandler.Login)
	api.Post("/register", authHandler.Register)

	// Profile
	api.Get("/profile", middlewares.Protected(), userHandler.GetProfile)

	// Games
	gameGroup := api.Group("/games", middlewares.Protected())
	gameGroup.Get("/", gameHandler.Index)
	gameGroup.Get("/current", gameHandler.GetCurrentGame)
	gameGroup.Get("/:code", gameHandler.Show)
	gameGroup.Post("/", gameHandler.Store)
	gameGroup.Delete("/leave", gameHandler.Leave)
	gameGroup.Post("/join/:id", gameHandler.Join)

	// Currencies
	api.Get("/currencies", currencyHandler.Index)

	// 404
	app.Use(func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound) // => 404 "Not Found"
	})
}
