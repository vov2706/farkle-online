package services

import (
	"app/database"
	"app/models"
	"app/repositories"
	"errors"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo     *repositories.UserRepository
	currencyRepo *repositories.CurrencyRepository
	balanceRepo  *repositories.BalanceRepository
}

func NewUserService(
	userRepo *repositories.UserRepository,
	currencyRepo *repositories.CurrencyRepository,
	balanceRepo *repositories.BalanceRepository,
) *UserService {
	return &UserService{
		userRepo:     userRepo,
		currencyRepo: currencyRepo,
		balanceRepo:  balanceRepo,
	}
}

func (service UserService) GetUserByUsername(username string) (*models.User, error) {
	return service.userRepo.GetUserByUsername(username)
}

func (service UserService) GetUserById(id uint) (*models.User, error) {
	return service.userRepo.GetUserById(id)
}

func (service UserService) CreateUser(username, password string) (*models.User, error) {
	currency, err := service.currencyRepo.FindBySlug(models.BRONZE)

	if err != nil {
		return nil, errors.New("currency not found")
	}

	var user *models.User

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		userRepo := service.userRepo.WithTx(tx)
		balanceRepo := service.balanceRepo.WithTx(tx)

		user, err = userRepo.CreateUser(username, password)

		if err != nil {
			return errors.New("failed to create user")
		}

		balance, err := balanceRepo.CreateBalance(user.ID, currency.ID, 1000)

		if err != nil {
			return errors.New("failed to create balance")
		}

		user.Balances = []models.Balance{*balance}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}
