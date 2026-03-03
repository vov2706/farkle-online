package models

import "time"

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Username    string    `json:"username" gorm:"uniqueIndex;not null;type:varchar(255)"`
	Password    string    `json:"password" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Balances    []Balance
	Friends     []*User `gorm:"many2many:user_friends"`
	Games       []Game  `gorm:"many2many:game_user"`
	CurrentGame *Game   `json:"current_game,omitempty" gorm:"-"`
}
