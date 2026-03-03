package database

import (
	"app/models"
	"time"

	"gorm.io/gorm/clause"
)

func SeedAll() {
	seedCurrencies()
}

type CurrencyItem struct {
	Slug string
	Name string
}

func seedCurrencies() {
	now := time.Now()
	items := []models.Currency{
		{Slug: models.BRONZE, Name: "Bronze", CreatedAt: now, UpdatedAt: now},
		// 		{Slug: models.SILVER, Name: "Silver", CreatedAt: now, UpdatedAt: now},
		// 		{Slug: models.GOLD, Name: "Gold", CreatedAt: now, UpdatedAt: now},
	}

	DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "slug"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "updated_at"}),
	}).Create(&items)
}
