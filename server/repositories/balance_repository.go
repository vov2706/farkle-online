package repositories

import (
	"app/models"
	"context"

	"gorm.io/gorm"
)

type BalanceRepository struct {
	DB *gorm.DB
}

func NewBalanceRepository(DB *gorm.DB) *BalanceRepository {
	return &BalanceRepository{DB: DB}
}

func (repo *BalanceRepository) WithTx(tx *gorm.DB) *BalanceRepository {
	repo.DB = tx

	return repo
}

func (repo *BalanceRepository) FindByUserAndCurrency(user models.User, currencyId uint) (*models.Balance, error) {
	if len(user.Balances) > 0 {
		for _, b := range user.Balances {
			if b.CurrencyID == currencyId {
				return &b, nil
			}
		}
	}

	ctx := context.Background()
	balance, err := gorm.G[models.Balance](repo.DB).
		Where("user_id = ?", user.ID).
		Where("currency_id = ?", currencyId).
		First(ctx)

	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (repo *BalanceRepository) CreateBalance(userId, currencyId, amount uint) (*models.Balance, error) {
	ctx := context.Background()
	balance := models.Balance{
		UserID:     userId,
		CurrencyID: currencyId,
		Amount:     amount,
	}
	err := gorm.G[models.Balance](repo.DB).Create(ctx, &balance)

	if err != nil {
		return nil, err
	}

	return &balance, nil
}
