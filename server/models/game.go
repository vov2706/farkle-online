package models

import "time"

type Game struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Code          string     `json:"code" gorm:"uniqueIndex; not null; type:varchar(255)"`
	CurrencyID    uint       `json:"currency_id" gorm:"index; not null"`
	CreatorID     uint       `json:"creator_id" gorm:"index; not null"`
	WinnerID      *uint      `json:"winner_id" gorm:"index;"`
	Bet           uint       `json:"bet" gorm:"not null"`
	WinningPoints uint       `json:"winning_points" gorm:"not null"`
	JoinType      string     `json:"join_type" gorm:"type:varchar(255); default:'anyone'; not null; index"`
	StartedAt     *time.Time `json:"started_at" gorm:"index"`
	FinishedAt    *time.Time `json:"finished_at" gorm:"index"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Currency      Currency   `json:"currency"`
	Creator       User       `json:"creator"`
	Winner        *User      `json:"winner"`
	Players       []User     `json:"users" gorm:"many2many:game_user;"`
	PlayersCount  uint       `json:"players_count"`
}

type GameUser struct {
	UserID uint `json:"user_id" gorm:"primaryKey; index; not null"`
	GameID uint `json:"game_id" gorm:"primaryKey; index; not null"`
}

func (GameUser) TableName() string {
	return "game_user"
}
