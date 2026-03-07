package responses

import (
	"app/models"
	"time"
)

type GameResource struct {
	Code          string            `json:"code"`
	Bet           uint              `json:"bet"`
	WinningPoints uint              `json:"winning_points"`
	JoinType      string            `json:"join_type"`
	StartedAt     *time.Time        `json:"started_at"`
	Currency      *CurrencyResource `json:"currency"`
	Players       []PlayerResource  `json:"players"`
	Creator       *UserResource     `json:"creator"`
	PlayersCount  uint              `json:"players_count"`
}

type PlayerResource struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	IsHost   bool   `json:"is_host"`
}

func NewPlayerResource(user models.User, isHost bool) PlayerResource {
	return PlayerResource{
		ID:       user.ID,
		Username: user.Username,
		IsHost:   isHost,
	}
}

func NewGameResource(game models.Game) GameResource {
	rs := GameResource{
		StartedAt:     game.StartedAt,
		Code:          game.Code,
		Bet:           game.Bet,
		WinningPoints: game.WinningPoints,
		JoinType:      game.JoinType,
		PlayersCount:  game.PlayersCount,
	}
	if game.Creator.ID != 0 {
		creator := NewUserResource(game.Creator)
		rs.Creator = &creator
	}

	if game.Currency.ID != 0 {
		c := NewCurrencyResource(game.Currency)
		rs.Currency = &c
	}

	if game.Players != nil {
		for _, p := range game.Players {
			rs.Players = append(rs.Players, NewPlayerResource(p, p.ID == game.CreatorID))
		}
	}

	return rs
}
