package container

import (
	"app/database"
	"app/http/handlers"
	"app/repositories"
	"app/services"
	"app/ws"

	"gorm.io/gorm"
)

type Container struct {
	DB       *gorm.DB
	Handlers *Handlers
	Services *Services
	Repos    *Repositories
	WS       *Websockets
}

type Websockets struct {
	Hub        *ws.Hub
	Router     *ws.Router
	Authorizer *ws.Authorizer
}

type Handlers struct {
	Auth     *handlers.AuthHandler
	Currency *handlers.CurrencyHandler
	Game     *handlers.GameHandler
	User     *handlers.UserHandler
}

type Services struct {
	Auth *services.AuthService
	Game *services.GameService
	User *services.UserService
}

type Repositories struct {
	Balance  *repositories.BalanceRepository
	Currency *repositories.CurrencyRepository
	Game     *repositories.GameRepository
	User     *repositories.UserRepository
}

func New() *Container {
	db := database.DB

	repos := &Repositories{
		Balance:  repositories.NewBalanceRepository(db),
		Game:     repositories.NewGameRepository(db),
		User:     repositories.NewUserRepository(db),
		Currency: repositories.NewCurrencyRepository(db),
	}

	svcs := &Services{
		Game: services.NewGameService(repos.Balance, repos.Currency, repos.Game),
		Auth: services.NewAuthService(repos.User),
		User: services.NewUserService(repos.User, repos.Currency, repos.Balance),
	}

	hub := ws.NewHub()
	websockets := &Websockets{
		Hub:        hub,
		Authorizer: ws.NewAuthorizer(),
		Router:     ws.NewRouter(hub),
	}

	hs := &Handlers{
		Game:     handlers.NewGameHandler(svcs.Game, svcs.Auth, hub),
		Auth:     handlers.NewAuthHandler(svcs.Auth, svcs.User),
		User:     handlers.NewUserHandler(svcs.Auth, svcs.Game),
		Currency: handlers.NewCurrencyHandler(repos.Currency),
	}

	return &Container{
		DB:       db,
		Repos:    repos,
		Services: svcs,
		Handlers: hs,
		WS:       websockets,
	}
}
