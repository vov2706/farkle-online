package repositories

import (
	"app/models"
	"app/utils"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) WithTx(tx *gorm.DB) *UserRepository {
	repo.DB = tx

	return repo
}

func (repo *UserRepository) GetUserById(id uint, relations ...string) (*models.User, error) {
	ctx := context.Background()
	query := gorm.G[models.User](repo.DB).Where("id = ?", id)

	if len(relations) > 0 {
		query = utils.ApplyRelations(query, relations)
	}

	user, err := query.First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}

		return nil, errors.New("failed to get user")
	}

	return &user, nil
}

func (repo *UserRepository) GetUserByUsername(username string, relations ...string) (*models.User, error) {
	ctx := context.Background()
	query := gorm.G[models.User](repo.DB).Where("username = ?", username)

	if len(relations) > 0 {
		query = utils.ApplyRelations(query, relations)
	}

	user, err := query.First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}

		return nil, errors.New("failed to get user")
	}

	return &user, nil
}

func (repo *UserRepository) CreateUser(username string, password string) (*models.User, error) {
	ctx := context.Background()
	user := models.User{
		Username: username,
		Password: password,
	}

	err := gorm.G[models.User](repo.DB).Create(ctx, &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
