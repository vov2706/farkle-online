package database

import (
	"app/config"
	"app/models"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() {
	var err error

	dbName := config.Config("DB_NAME")

	DB, err = gorm.Open(sqlite.Open(dbName+".db"), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}

	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		panic(err)
	}

	fmt.Println("Connected to database")

	migrateAll()
	setupRelations()
	SeedAll()
}

func migrateAll() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Currency{},
		&models.Balance{},
		&models.Game{},
		&models.GameUser{},
	)

	if err != nil {
		panic(fmt.Errorf("failed to migrate database: %w", err))
	}

	fmt.Println("Database migrated successfully")
}

func setupRelations() {
	err := DB.SetupJoinTable(&models.Game{}, "Players", &models.GameUser{})

	if err != nil {
		panic(fmt.Errorf("failed to setup relations: %w", err))
	}

	err = DB.SetupJoinTable(&models.User{}, "Games", &models.GameUser{})

	if err != nil {
		panic(fmt.Errorf("failed to setup relations: %w", err))
	}
}
