package services

import (
	"app/http/inputs"
	"app/http/responses"
	"app/models"
	"app/repositories"
	"crypto/rand"
	"errors"
	"math/big"
)

type GameService struct {
	balanceRepo  *repositories.BalanceRepository
	currencyRepo *repositories.CurrencyRepository
	gameRepo     *repositories.GameRepository
}

func NewGameService(
	balanceRepo *repositories.BalanceRepository,
	currencyRepo *repositories.CurrencyRepository,
	gameRepo *repositories.GameRepository,
) *GameService {
	return &GameService{
		balanceRepo:  balanceRepo,
		currencyRepo: currencyRepo,
		gameRepo:     gameRepo,
	}
}

func (service *GameService) PaginatePublicGames(authUserId, page, perPage uint, search string) (*responses.PaginatedResponse, error) {
	games, total, err := service.gameRepo.GetNewGames(authUserId, page, perPage, search)

	if err != nil {
		return nil, err
	}

	results := make([]responses.GameResource, 0)

	if len(games) > 0 {
		for _, game := range games {
			results = append(results, responses.NewGameResource(game))
		}
	}

	paginationMeta := responses.NewPaginationMeta(page, perPage, total)
	res := responses.NewPaginatedResponse(results, paginationMeta)

	return &res, nil
}

func (service *GameService) GetGameByCode(code string) *models.Game {
	game, err := service.gameRepo.FindGameByCode(code, "Players", "Currency")

	if err != nil {
		return nil
	}

	return &game
}

func (service *GameService) GetGameById(id uint) *models.Game {
	game, err := service.gameRepo.FindGameById(id, "Players", "Currency")

	if err != nil {
		return nil
	}

	return game
}

func (service *GameService) CreateGame(authUser *models.User, input *inputs.CreateGameInput) (*models.Game, error) {
	currency, err := service.currencyRepo.FindById(input.CurrencyID)

	if err != nil {
		return nil, errors.New("failed to get currency")
	}

	userBalance, err := service.balanceRepo.FindByUserAndCurrency(*authUser, currency.ID)

	if err != nil {
		return nil, errors.New("failed to get user balance")
	}

	if userBalance.Amount < input.Bet {
		return nil, errors.New("insufficient funds")
	}

	code, err := service.generateGameCode(6)

	game, err := service.gameRepo.CreateGame(*authUser, *input, code)

	if err != nil {
		return nil, errors.New("failed to create game")

	}

	game, err = service.gameRepo.FindGameById(game.ID, "Players")

	if err != nil {
		return nil, errors.New("not found game")
	}

	game.Currency = currency

	return game, nil
}

func (service *GameService) JoinToGame(authUser *models.User, gameId uint) (*models.Game, error) {
	game, err := service.gameRepo.FindGameById(gameId)

	if err != nil {
		return nil, err
	}

	if err = service.gameRepo.AttachUserToGame(*authUser, *game); err != nil {
		return nil, err
	}

	return game, nil
}

func (service *GameService) GetCurrentGame(authUser models.User) *models.Game {
	game, err := service.gameRepo.FindCurrentGame(authUser.ID)

	if err != nil {
		return nil
	}

	return &game
}

func (service *GameService) LeaveCurrentGame(authUser models.User) error {
	game := service.GetCurrentGame(authUser)
	var err error

	if game != nil {
		if authUser.ID == game.CreatorID {
			err = service.gameRepo.DeleteGame(authUser, *game)
		} else {
			err = service.gameRepo.DetachUserFromGame(authUser, *game)
		}
	}

	return err
}

func (service *GameService) generateGameCode(length int) (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, length)

	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}

	return string(b), nil
}
