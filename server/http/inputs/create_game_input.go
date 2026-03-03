package inputs

import (
	"app/database"
	"errors"
	"fmt"
)

const (
	Anyone      = "anyone"
	OnlyFriends = "friends"
	ByLink      = "link"
)

const WinningPointsLimit = 20000
const WinningPointsMinimum = 3000

type CreateGameInput struct {
	CurrencyID    uint   `json:"currency_id" validate:"required,int,gt:0"`
	Bet           uint   `json:"bet" validate:"required,int,gt:0"`
	WinningPoints uint   `json:"winning_points" validate:"required,int,gt:0"`
	JoinType      string `json:"join_type" validate:"required,string"`
}

func (input CreateGameInput) Validate() error {
	if !input.isCurrencyIsValid() {
		return errors.New("invalid currency")
	}

	if !input.isValidJoinType() {
		return fmt.Errorf("invalid join type %s", input.JoinType)
	}

	if !input.lteWinningPointsLimit() {
		return fmt.Errorf("winning points limit exceeded %d points", WinningPointsLimit)
	}

	if !input.gteWinningPointsMinimum() {
		return fmt.Errorf("winning points must be at least %d points", WinningPointsMinimum)
	}

	return nil
}

func (input CreateGameInput) isCurrencyIsValid() bool {
	var exists bool

	err := database.DB.Raw(
		"SELECT EXISTS (SELECT 1 FROM currencies WHERE id = ?)",
		input.CurrencyID,
	).Scan(&exists).Error

	if err != nil {
		return false
	}

	return exists
}

func (input CreateGameInput) isValidJoinType() bool {
	switch input.JoinType {
	case Anyone, OnlyFriends, ByLink:
		return true
	default:
		return false
	}
}

func (input CreateGameInput) lteWinningPointsLimit() bool {
	return input.WinningPoints <= WinningPointsLimit
}

func (input CreateGameInput) gteWinningPointsMinimum() bool {
	return input.WinningPoints >= WinningPointsMinimum
}
