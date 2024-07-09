package dto

type GetAccountResponse struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

type GetAllAccountResponse struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}
