package repositories

import (
	"app/database"
	"app/models"
	"context"

	"gorm.io/gorm"
)

type CurrencyRepository struct {
	db *gorm.DB
}

func NewCurrencyRepository(db *gorm.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

func (repo CurrencyRepository) GetCurrencies() ([]models.Currency, error) {
	ctx := context.Background()

	return gorm.G[models.Currency](database.DB).Find(ctx)
}

func (repo CurrencyRepository) FindById(id uint) (models.Currency, error) {
	ctx := context.Background()

	return gorm.G[models.Currency](database.DB).
		Where("id = ?", id).
		First(ctx)
}

func (repo CurrencyRepository) FindBySlug(slug string) (models.Currency, error) {
	ctx := context.Background()

	return gorm.G[models.Currency](database.DB).
		Where("slug = ?", slug).
		First(ctx)
}
