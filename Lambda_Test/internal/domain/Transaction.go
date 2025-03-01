package domain

type Transaction struct {
	Id     string  `json:"transaction_id"`
	Amount float64 `json:"amount"`
}
