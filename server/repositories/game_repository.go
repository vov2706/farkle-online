package repositories

import (
	"app/database"
	"app/http/inputs"
	"app/models"
	"app/utils"
	"context"
	"strings"

	"gorm.io/gorm"
)

type GameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{db: db}
}

func (repo *GameRepository) CreateGame(user models.User, input inputs.CreateGameInput, code string) (*models.Game, error) {
	var game models.Game
	ctx := context.Background()

	err := database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		game := models.Game{
			Code:          code,
			CreatorID:     user.ID,
			CurrencyID:    input.CurrencyID,
			Bet:           input.Bet,
			WinningPoints: input.WinningPoints,
			JoinType:      input.JoinType,
		}

		if err := tx.Create(&game).Error; err != nil {
			return err
		}

		if err := tx.Model(&user).Association("Games").Append(game); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (repo *GameRepository) FindGameById(id uint, relations ...string) (*models.Game, error) {
	ctx := context.Background()
	query := gorm.G[models.Game](database.DB).Where("id = ?", id)

	if len(relations) > 0 {
		query = utils.ApplyRelations(query, relations)
	}

	game, err := query.Select("*").First(ctx)

	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (repo *GameRepository) FindGameByCode(code string, relations ...string) (models.Game, error) {
	ctx := context.Background()
	query := gorm.G[models.Game](database.DB).Where("code = ?", code)

	if len(relations) > 0 {
		query = utils.ApplyRelations(query, relations)
	}

	return query.First(ctx)
}

func (repo *GameRepository) GetNewGames(authUserId, page, perPage uint, search string) ([]models.Game, uint, error) {
	var games []models.Game
	var total int64

	query := database.DB.Where("started_at IS NULL").Where("creator_id != ?", authUserId)
	trimmedSearch := strings.TrimSpace(search)

	if trimmedSearch != "" {
		query = query.Where("code LIKE ?", "%"+trimmedSearch+"%")
	}

	if err := query.Model(&models.Game{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query.Scopes(database.Paginate(page, perPage), database.OrderByDesc("id")).
		Preload("Currency", nil).
		Preload("Creator", nil).
		Select("*", "(select count(*) from game_user where game_id = games.id) as players_count").
		Find(&games)

	return games, uint(total), nil
}

func (repo *GameRepository) FindCurrentGame(userId uint) (models.Game, error) {
	ctx := context.Background()
	game, err := gorm.G[models.Game](database.DB).
		Where("started_at IS NULL").
		Where(`
			EXISTS (
				SELECT 1
				FROM game_user
				WHERE game_user.game_id = games.id
				AND game_user.user_id = ?
			)
		`, userId).
		First(ctx)

	return game, err
}

func (repo *GameRepository) DeleteGame(creator models.User, game models.Game) error {
	ctx := context.Background()

	return database.DB.Transaction(func(tx *gorm.DB) error {
		err := database.DB.Model(&game).Association("Players").Clear()

		if err != nil {
			return err
		}

		_, err = gorm.G[models.Game](database.DB).
			Where("creator_id = ?", creator.ID).
			Where("id = ?", game.ID).
			Where("started_at IS NULL").
			Delete(ctx)

		if err != nil {
			return err
		}

		return nil
	})
}

func (repo *GameRepository) DetachUserFromGame(user models.User, game models.Game) error {
	return database.DB.Model(&user).Association("Games").Delete(game)
}

func (repo *GameRepository) AttachUserToGame(user models.User, game models.Game) error {
	return database.DB.Model(&user).Association("Games").Append(game)
}
