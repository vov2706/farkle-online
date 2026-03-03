package responses

type BalanceResource struct {
	Amount   uint             `json:"amount"`
	Currency CurrencyResource `json:"currency"`
}

func NewBalanceResource(amount uint, currency CurrencyResource) BalanceResource {
	return BalanceResource{Amount: amount, Currency: currency}
}
