package responses

import "app/models"

type UserResponse struct {
	Data UserResource `json:"data"`
}

type UserResource struct {
	ID          uint              `json:"id"`
	Username    string            `json:"username"`
	CurrentGame *models.Game      `json:"current_game,omitempty"`
	Balances    []BalanceResource `json:"balances,omitempty"`
}

func NewUserResource(user models.User) UserResource {
	var balances []BalanceResource

	for _, balance := range user.Balances {
		balances = append(balances, NewBalanceResource(
			balance.Amount,
			NewCurrencyResource(balance.Currency),
		))
	}

	return UserResource{
		ID:          user.ID,
		Username:    user.Username,
		Balances:    balances,
		CurrentGame: user.CurrentGame,
	}
}
