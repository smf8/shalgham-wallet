package handler

type TransactionRequest struct {
	PhoneNumber string  `json:"phone_number"`
	Amount      float64 `json:"amount"`
}

type ProfileCreateRequest struct {
	PhoneNumber string  `json:"phone_number"`
	Balance     float64 `json:"balance"`
}
